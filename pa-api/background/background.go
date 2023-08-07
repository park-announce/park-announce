package background

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"os"

	"github.com/go-redis/redis"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/global"
	"github.com/park-announce/pa-api/server"
	"github.com/park-announce/pa-api/service"
	kafka "github.com/segmentio/kafka-go"
)

func NewBackgroundOperation() *BackgroundOperation {
	return &BackgroundOperation{}
}

type BackgroundOperation struct {
}

func (backgroundOperation *BackgroundOperation) GetGlobalInstanceId(wgMain *sync.WaitGroup, redis *redis.Client, redisLockObtained chan bool, redisLockQuit chan os.Signal, redisLockName chan string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
		}
	}()
	defer func() {
		wgMain.Done()
	}()
	var lockName string
	var lockValue string
	for {

		lockName = fmt.Sprintf("global-api-instance-id-%d", global.GetInstanceId())

		set, setErr := redis.SetNX(lockName, lockValue, 0).Result()

		if setErr != nil {
			log.Println(setErr)
			global.IncrementInstanceId()
			continue
		}
		if !set {
			log.Println("value could not be set")
			global.IncrementInstanceId()
			continue
		}

		break
	}

	redisLockName <- lockName
	redisLockObtained <- true
	<-redisLockQuit

	log.Println("interrupt detect in redis go rouitine")
	lockValueFromRedis, getErr := redis.Get(lockName).Result()

	log.Printf("lockName : %s, lockValue : %s, lockValueFromRedis : %s", lockName, lockValue, lockValueFromRedis)

	if getErr != nil {
		log.Println(getErr)
	}

	if lockValueFromRedis == lockValue {
		deleteResult, delErr := redis.Del(lockName).Result()

		if delErr != nil {
			log.Println(delErr)
		}
		if deleteResult > 0 {
			log.Printf("redis key delete result. key : %s, value : %s, result : %d", lockName, lockValue, deleteResult)
		}

	}
}

func (backgroundOperation *BackgroundOperation) RedisHeartBeat(wgMain *sync.WaitGroup, redis *redis.Client, redisHeartBeatQuit chan os.Signal, lockName string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
		}
	}()
	defer func() {
		wgMain.Done()
	}()

	for {

		select {

		case <-redisHeartBeatQuit:

			return

		default:

			_, setErr := redis.Expire(lockName, time.Second*30).Result()

			if setErr != nil {
				log.Println(setErr)
				continue
			}

			time.Sleep(time.Second)

		}

	}
}

func (backgroundOperation *BackgroundOperation) StartServer(wgMain *sync.WaitGroup, hub *service.SocketHub, httpServerQuit chan os.Signal) {
	defer func() {
		wgMain.Done()
	}()

	go func() {

		defer func() {
			if err := recover(); err != nil {

				log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			}
		}()

		server.NewServer(hub).Run(":8000")
	}()

	<-httpServerQuit
	log.Println("interrupt detect in http go rouitine")

}

func (backgroundOperation *BackgroundOperation) ConsumeClientResponse(wgMain *sync.WaitGroup, redis *redis.Client, hub *service.SocketHub, kafkaConsumerQuit chan os.Signal) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
		}
	}()

	defer func() {
		wgMain.Done()
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     fmt.Sprintf("client_response_%d", global.GetInstanceId()),
		GroupID:   "pa-api-consumer-group",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		Topic:                  "dead_letter_messages",
		AllowAutoTopicCreation: true,
	}

	defer func() {

		if err := writer.Close(); err != nil {
			log.Println("failed to close writer:", err)
		}
	}()

	for {

		select {

		case <-kafkaConsumerQuit:
			if err := reader.Close(); err != nil {
				log.Println("failed to close reader:", err)
			}
			return
		default:
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("error :", err)
				break
			}

			fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

			var clientKafkaResponseMessage entity.ClientKafkaResponseMessage

			err = json.Unmarshal(m.Value, &clientKafkaResponseMessage)
			if err != nil {
				log.Println("error :", err)
				break
			}

			var clientSocketResponseMessage entity.ClientSocketResponseMessage
			err = json.Unmarshal(m.Value, &clientSocketResponseMessage)
			if err != nil {
				log.Println("error :", err)
				break
			}

			message, err := json.Marshal(clientSocketResponseMessage)

			if err != nil {
				log.Println("error :", err)
				break
			}

			if !hub.SendMessageIfClientExist(clientKafkaResponseMessage.ClientId, message) {
				//produce this message to dead_letter_message topic
				log.Printf("socket connection not found with id : %s, message is sending to dead_letter_messages topic", clientKafkaResponseMessage.ClientId)
				messages := []kafka.Message{
					{
						Value: m.Value,
					},
				}

				const retries = 3
				for i := 0; i < retries; i++ {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					// attempt to create topic prior to publishing the message
					err = writer.WriteMessages(ctx, messages...)
					if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
						time.Sleep(time.Millisecond * 250)
						continue
					}

					if err != nil {
						log.Println("error :", err)
					}
					break
				}
				break
			} else {
				transactionUniqueKey := fmt.Sprintf("transaction|%s|%s|%s", clientKafkaResponseMessage.ClientId, clientSocketResponseMessage.Operation, clientSocketResponseMessage.TransactionId)

				_, err := redis.Del(transactionUniqueKey).Result()
				if err != nil {
					log.Println("error :", err)
					return
				}
			}

			reader.CommitMessages(context.Background(), m)
		}

	}

}

func (backgroundOperation *BackgroundOperation) ConsumeDeadLetterMessages(wgMain *sync.WaitGroup, redis *redis.Client, hub *service.SocketHub, kafkaConsumerQuit chan os.Signal) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error in recover : %v, stack : %s", err, string(debug.Stack()))
		}
	}()

	defer func() {
		wgMain.Done()
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "dead_letter_messages",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	for {

		select {

		case <-kafkaConsumerQuit:
			if err := reader.Close(); err != nil {
				log.Println("failed to close reader:", err)
			}
			return
		default:
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("error :", err)
				break
			}
			fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

			var clientKafkaResponseMessage entity.ClientKafkaResponseMessage

			err = json.Unmarshal(m.Value, &clientKafkaResponseMessage)
			if err != nil {
				log.Println("error :", err)
				break
			}

			var clientSocketResponseMessage entity.ClientSocketResponseMessage
			err = json.Unmarshal(m.Value, &clientSocketResponseMessage)
			if err != nil {
				log.Println("error :", err)
				break
			}

			message, err := json.Marshal(clientSocketResponseMessage)

			if err != nil {
				log.Println("error :", err)
				break
			}

			result := hub.SendMessageIfClientExist(clientKafkaResponseMessage.ClientId, message)

			if result {
				transactionUniqueKey := fmt.Sprintf("transaction|%s|%s|%s", clientKafkaResponseMessage.ClientId, clientSocketResponseMessage.Operation, clientSocketResponseMessage.TransactionId)

				_, err := redis.Del(transactionUniqueKey).Result()
				if err != nil {
					log.Println("error :", err)
					return
				}
			}

			reader.CommitMessages(context.Background(), m)
		}

	}
}

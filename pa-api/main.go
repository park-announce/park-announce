package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
	kafka "github.com/segmentio/kafka-go"

	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/global"
	"github.com/park-announce/pa-api/server"
	"github.com/park-announce/pa-api/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var lockName string
var lockValue string = uuid.NewString()

func main() {

	httpServerQuit := make(chan os.Signal, 1)
	signal.Notify(httpServerQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	redisLockQuit := make(chan os.Signal, 1)
	signal.Notify(redisLockQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	redisHeartBeatQuit := make(chan os.Signal, 1)
	signal.Notify(redisHeartBeatQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	kafkaConsumerQuit := make(chan os.Signal, 1)
	signal.Notify(kafkaConsumerQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	redisLockObtained := make(chan bool)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	hub := service.NewSocketHub()

	wgMain := &sync.WaitGroup{}
	wgMain.Add(4)

	//try to get global instance id from redis using distributed lock.
	go func(_wgMain *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in redis go rouitine", r)
			}
		}()
		defer func() {
			_wgMain.Done()
		}()

		for {

			lockName = fmt.Sprintf("global-api-instance-id-%d", global.GetInstanceId())

			set, setErr := rdb.SetNX(lockName, lockValue, 0).Result()

			if setErr != nil {
				log.Fatal(setErr)
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

		redisLockObtained <- true
		<-redisLockQuit

		log.Println("interrupt detect in redis go rouitine")
		lockValueFromRedis, getErr := rdb.Get(lockName).Result()

		log.Printf("lockName : %s, lockValue : %s, lockValueFromRedis : %s", lockName, lockValue, lockValueFromRedis)

		if getErr != nil {
			log.Fatal(getErr)
		}

		if lockValueFromRedis == lockValue {
			deleteResult, delErr := rdb.Del(lockName).Result()

			if delErr != nil {
				log.Fatal(delErr)
			}
			if deleteResult > 0 {
				log.Printf("redis key delete result. key : %s, value : %s, result : %d", lockName, lockValue, deleteResult)
			}

		}

	}(wgMain)

	//main routine continious running after redis lock obtained
	<-redisLockObtained

	//send heartbeat message to redis in spesific interval setting expire value for distributed lock key in redis
	go func(_wgMain *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in redis heartbeat goroutine", r)
			}
		}()
		defer func() {
			_wgMain.Done()
		}()

		for {

			select {

			case <-redisHeartBeatQuit:

				return

			default:

				set, setErr := rdb.Expire(lockName, time.Second*30).Result()

				if setErr != nil {
					log.Fatal(setErr)
					continue
				}
				if !set {
					log.Printf("expire could not be set. lockname : %s", lockName)
					continue
				}

				time.Sleep(time.Second)

			}

		}

	}(wgMain)

	//handle socket message connection and message request.
	go func(_wgMain *sync.WaitGroup) {
		defer func() {
			_wgMain.Done()
		}()

		go func() {

			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in http go routine", r)
				}
			}()

			server.NewServer(hub).Run(":8000")
		}()

		<-httpServerQuit
		log.Println("interrupt detect in http go rouitine")

	}(wgMain)

	//consume message from client_response_{{instanceId}} topic and try to send message to related socket client
	//if there is no socket connection, sends this message to dead_letter_message topic
	go func(_wgMain *sync.WaitGroup) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in http go routine", r)
			}
		}()

		defer func() {
			_wgMain.Done()
		}()

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"kafka:9092"},
			Topic:     fmt.Sprintf("client_response_%d", global.GetInstanceId()),
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
				log.Fatal("failed to close writer:", err)
			}
		}()

		for {

			select {

			case <-kafkaConsumerQuit:
				if err := reader.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
				return
			default:
				m, err := reader.ReadMessage(context.Background())
				if err != nil {
					break
				}

				fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

				var clientKafkaResponseMessage entity.ClientKafkaResponseMessage

				err = json.Unmarshal(m.Value, &clientKafkaResponseMessage)
				if err != nil {
					break
				}

				if !hub.SendMessageIfClientExist(clientKafkaResponseMessage.ClientId, m.Value) {
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
							log.Fatalf("unexpected error %v", err)
						}
						break
					}
					break
				}

				reader.CommitMessages(context.Background(), m)
			}

		}

	}(wgMain)

	//consume message from dead_letter_messages topic and try to send message to related socket client
	go func(_wgMain *sync.WaitGroup) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in http go routine", r)
			}
		}()

		defer func() {
			_wgMain.Done()
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
					log.Fatal("failed to close reader:", err)
				}
				return
			default:
				m, err := reader.ReadMessage(context.Background())
				if err != nil {
					break
				}
				fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

				var clientKafkaResponseMessage entity.ClientKafkaResponseMessage

				err = json.Unmarshal(m.Value, &clientKafkaResponseMessage)
				if err != nil {
					log.Printf("unexpected error %v", err)
					break
				}

				hub.SendMessageIfClientExist(clientKafkaResponseMessage.ClientId, m.Value)

				reader.CommitMessages(context.Background(), m)
			}

		}

	}(wgMain)

	wgMain.Wait()
}

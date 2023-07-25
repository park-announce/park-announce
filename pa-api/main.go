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
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
	kafka "github.com/segmentio/kafka-go"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SocketMessage struct {
	ClientId      string      `json:"client_id"`
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	Data          interface{} `json:"data"`
}

type ClientKafkaResponseMessage struct {
	ClientId string      `json:"client_id"`
	ApiId    string      `json:"api_id"`
	Data     interface{} `json:"data"`
}

type ClientKafkaRequestMessage struct {
	ClientId      string      `json:"client_id"`
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	ApiId         string      `json:"api_id"`
	Data          interface{} `json:"data"`
}

type SendSocketMessage struct {
	ClientId string      `json:"client_id"`
	Data     interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var instanceId int = 0
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

			lockName = fmt.Sprintf("global-api-instance-id-%d", instanceId)

			set, setErr := rdb.SetNX(lockName, lockValue, 0).Result()

			if setErr != nil {
				log.Fatal(setErr)
				instanceId++
				continue
			}
			if !set {
				log.Println("value could not be set")
				instanceId++
				continue
			}

			break
		}

		log.Printf("instanceId : %d, lockValue : %s", instanceId, lockValue)
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
					log.Println("expire could not be set")
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
			// Initialize router
			router := mux.NewRouter()

			// Define API endpoints
			router.HandleFunc("/socket/connect", connectSocket)
			router.HandleFunc("/socket/messages", sendSampleMessage)

			log.Println("Server started on port 8000")
			log.Fatal(http.ListenAndServe(":8000", router))
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

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"kafka:9092"},
			Topic:     fmt.Sprintf("client_response_%d", instanceId),
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
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
				return
			default:
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					break
				}

				fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

				var clientKafkaResponseMessage ClientKafkaResponseMessage

				err = json.Unmarshal(m.Value, &clientKafkaResponseMessage)
				if err != nil {
					break
				}

				conn := connections[clientKafkaResponseMessage.ClientId]

				if conn == nil {
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

				// data, err := json.Marshal(clientKafkaResponseMessage.Data)
				// if err != nil {
				// 	log.Printf("unexpected error %v", err)
				// 	break
				// }

				// write message to related socket client
				if err := conn.WriteMessage(websocket.TextMessage, m.Value); err != nil {
					return
				}
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

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"kafka:9092"},
			Topic:     "dead_letter_messages",
			Partition: 0,
			MaxBytes:  10e6, // 10MB
		})

		for {

			select {

			case <-kafkaConsumerQuit:
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
				return
			default:
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					break
				}
				fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

				var clientKafkaResponseMessage ClientKafkaResponseMessage

				err = json.Unmarshal(m.Value, &clientKafkaResponseMessage)
				if err != nil {
					log.Printf("unexpected error %v", err)
					break
				}

				conn := connections[clientKafkaResponseMessage.ClientId]

				if conn == nil {
					break
				}

				// data, err := json.Marshal(clientKafkaResponseMessage.Data)
				// if err != nil {
				// 	log.Printf("unexpected error %v", err)
				// 	break
				// }

				// write message to related socket client
				if err := conn.WriteMessage(websocket.TextMessage, m.Value); err != nil {
					return
				}
			}

		}

	}(wgMain)

	wgMain.Wait()
}

func sendSampleMessage(w http.ResponseWriter, r *http.Request) {

	var sendSocketMessage SendSocketMessage
	json.NewDecoder(r.Body).Decode(&sendSocketMessage)

	conn := connections[sendSocketMessage.ClientId]

	data, err := json.Marshal(sendSocketMessage.Data)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return
	}

	// Write message back to browser
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("unexpected error %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("OK")
}

func connectSocket(w http.ResponseWriter, r *http.Request) {

	conn, connErr := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	if connErr != nil {
		log.Println(connErr)
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		Topic:                  "client_request",
		AllowAutoTopicCreation: true,
	}

	defer func() {

		if err := writer.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var socketMessage SocketMessage
		deSerializationError := json.Unmarshal(msg, &socketMessage)
		if deSerializationError != nil {
			log.Println(deSerializationError)
			continue
		}

		//validate ClientId

		if socketMessage.ClientId == "" {
			log.Println("invalid client id")
			return
		}

		connections[socketMessage.ClientId] = conn

		clientKafkaRequestMessage := &ClientKafkaRequestMessage{
			ClientId:      socketMessage.ClientId,
			ApiId:         fmt.Sprintf("%d", instanceId),
			TransactionId: socketMessage.TransactionId,
			Operation:     socketMessage.Operation,
			Data:          socketMessage.Data,
		}

		data, err := json.Marshal(clientKafkaRequestMessage)

		if err != nil {
			return
		}

		messages := []kafka.Message{
			{
				Value: data,
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

	}
}

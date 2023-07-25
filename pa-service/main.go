package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

type ClientKafkaRequestMessage struct {
	ClientId      string      `json:"client_id"`
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	ApiId         string      `json:"api_id"`
	Data          interface{} `json:"data"`
}

type ClientKafkaResponseMessage struct {
	ClientId      string      `json:"client_id"`
	ApiId         string      `json:"api_id"`
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	Data          interface{} `json:"data"`
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type FindNearbyLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
}

type CreateParkLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type IOperation interface {
	Do(data interface{}) (error, interface{})
}

type FindLocationsNearbyOperation struct{}

func (o *FindLocationsNearbyOperation) Do(data interface{}) (error, interface{}) {

	clientKafkaRequestMessage := data.(ClientKafkaRequestMessage)

	findNearbyLocationRequestData, err := json.Marshal(clientKafkaRequestMessage.Data)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	var findNearbyLocationRequest FindNearbyLocationRequest

	err = json.Unmarshal(findNearbyLocationRequestData, &findNearbyLocationRequest)
	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	rows, err := db.Query("SELECT ST_X(geog::geometry) as longitude, ST_Y(geog::geometry) as latitude FROM pa_locations WHERE status = $1 and ST_DWithin(geog, ST_MakePoint($2, $3)::geography, $4)", 0, findNearbyLocationRequest.Longitude, findNearbyLocationRequest.Latitude, findNearbyLocationRequest.Distance)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		if err := rows.Scan(&location.Longitude, &location.Latitude); err != nil {
			log.Printf("unexpected error %v", err)
		}
		locations = append(locations, location)
	}
	return nil, locations
}

type CreateParkLocationOperation struct{}

func (o *CreateParkLocationOperation) Do(data interface{}) (error, interface{}) {

	clientKafkaRequestMessage := data.(ClientKafkaRequestMessage)

	createParkLocationRequestData, err := json.Marshal(clientKafkaRequestMessage.Data)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	var createParkLocationRequest CreateParkLocationRequest

	err = json.Unmarshal(createParkLocationRequestData, &createParkLocationRequest)
	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	result, err := db.Exec("INSERT INTO pa_locations (id,geog,owner_id) VALUES($1,ST_MakePoint($2, $3)::geography,$4)", uuid.New().String(), createParkLocationRequest.Longitude, createParkLocationRequest.Latitude, clientKafkaRequestMessage.ClientId)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	count, err := result.RowsAffected()

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	return err, count

}

var operations map[string]IOperation = make(map[string]IOperation, 0)

var db *sql.DB

func main() {

	operations["get_locations_nearby"] = &FindLocationsNearbyOperation{}
	operations["create_park_location"] = &CreateParkLocationOperation{}

	//initialize postgres client
	connStr := "postgres://park_announce:PosgresDb1591*@db/padb?sslmode=disable"
	var dbError error
	db, dbError = sql.Open("postgres", connStr)
	if dbError != nil {
		log.Println(dbError)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("failed to close db:", err)
		}
	}()

	kafkaConsumerQuit := make(chan os.Signal, 1)
	signal.Notify(kafkaConsumerQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	wgMain := &sync.WaitGroup{}
	wgMain.Add(1)

	//consume message from client_request topic and do sth and writes result to client_response_* topic
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
			Topic:     "client_request",
			GroupID:   "pa-service-consumer-group",
			Partition: 0,
			MaxBytes:  10e6, // 10MB
		})

		//initialize kafka producer
		writer := &kafka.Writer{
			Addr:                   kafka.TCP("kafka:9092"),
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
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
				return
			default:
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					log.Printf("unexpected error %v", err)
					break
				}

				fmt.Printf("message consumed from topic : %s,  offset : %d, key : %s, value : %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))

				var clientKafkaRequestMessage ClientKafkaRequestMessage

				err = json.Unmarshal(m.Value, &clientKafkaRequestMessage)
				if err != nil {
					log.Printf("unexpected error %v", err)
					break
				}

				operation := operations[clientKafkaRequestMessage.Operation]

				if operation == nil {
					log.Printf("unimplemented operation : %s", clientKafkaRequestMessage.Operation)
					break
				}

				err, responseData := operation.Do(clientKafkaRequestMessage)

				if err != nil {
					log.Printf("unexpected error %v", err)
					break
				}

				clientKafkaResponseMessage := &ClientKafkaResponseMessage{
					ClientId:      clientKafkaRequestMessage.ClientId,
					ApiId:         clientKafkaRequestMessage.ApiId,
					TransactionId: clientKafkaRequestMessage.TransactionId,
					Operation:     clientKafkaRequestMessage.Operation,
					Data:          responseData,
				}

				clientKafkaResponseMessageData, err := json.Marshal(clientKafkaResponseMessage)

				if err != nil {
					log.Printf("unexpected error %v", err)
				}

				messages := []kafka.Message{
					{
						Topic: fmt.Sprintf("client_response_%s", clientKafkaRequestMessage.ApiId),
						Value: clientKafkaResponseMessageData,
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
						log.Printf("unexpected error %v", err)
					}

					fmt.Printf("message produced to topic : %s, value : %s\n", messages[0].Topic, string(messages[0].Value))

					break
				}
			}

		}

	}(wgMain)

	wgMain.Wait()
}

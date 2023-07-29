package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-redis/redis"
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
	IsSuccess     bool        `json:"is_success"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationWithDistance struct {
	Id         string  `json:"id"`
	DistanceTo float64 `json:"distance_to"`
	Location
}

type LocationEntity struct {
	Id             string         `json:"id"`
	OwnerId        string         `json:"owner_id"`
	AssignedUserId sql.NullString `json:"assigned_user_id,omitempty"`
	Status         int16          `json:"status"`
	Location
}

type FindNearbyLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
	Count     int32   `json:"count"`
}

type FindNearbyLocationResponse struct {
	Duration  int32                  `json:"duration"`
	Locations []LocationWithDistance `json:"locations"`
}

type CreateParkLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ReserveParkLocationRequest struct {
	Id string `json:"id"`
}

type ScheduleParkLocationAvailabiltyRequest struct {
	Id            string       `json:"id"`
	Latitude      float64      `json:"latitude"`
	Longitude     float64      `json:"longitude"`
	ScheduleType  ScheduleType `json:"schedule_type"`
	ScheduledTime int64        `json:"scheduled_time"`
}

type ScheduleType int16

const (
	ByLocation ScheduleType = 0
	ById       ScheduleType = 1
)

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

	findNearbyLocationRequest.Count = int32(math.Min(float64(findNearbyLocationRequest.Count), float64(maxSearchResult)))

	findNearbyLocationRequest.Distance = math.Min(float64(findNearbyLocationRequest.Distance), float64(maxDistance))

	rows, err := db.Query("SELECT id, ST_X(geog::geometry) as longitude, ST_Y(geog::geometry) as latitude, ST_Distance(geog,ST_MakePoint($2, $3)::geography) as distance_to FROM pa_locations WHERE status = $1 and ST_DWithin(geog, ST_MakePoint($2, $3)::geography, $4) order by distance_to asc limit $5", 0, findNearbyLocationRequest.Longitude, findNearbyLocationRequest.Latitude, findNearbyLocationRequest.Distance, findNearbyLocationRequest.Count)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	defer rows.Close()

	var locations []LocationWithDistance
	for rows.Next() {
		var location LocationWithDistance
		if err := rows.Scan(&location.Id, &location.Longitude, &location.Latitude, &location.DistanceTo); err != nil {
			log.Printf("unexpected error %v", err)
		}
		locations = append(locations, location)
		/*
			//check is this location checked out for reservervation for a spesific duration
			result, _ := rdb.SetNX(fmt.Sprintf("reservation-lock-%s", location.Id), clientKafkaRequestMessage.ClientId, reservationLockDurationTime).Result()

			if result {
				//if not checkout for rezervation, add to list
				locations = append(locations, location)
			}

		*/

	}

	return nil, &FindNearbyLocationResponse{Duration: reservationLockDurationSeconds, Locations: locations}
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

type ReserveParkLocationOperation struct{}

func (o *ReserveParkLocationOperation) Do(data interface{}) (error, interface{}) {

	clientKafkaRequestMessage := data.(ClientKafkaRequestMessage)

	reserveParkLocationRequestData, err := json.Marshal(clientKafkaRequestMessage.Data)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	var reserveParkLocationRequest ReserveParkLocationRequest

	err = json.Unmarshal(reserveParkLocationRequestData, &reserveParkLocationRequest)
	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	/*
		//check is location reserved for this user
		reservedUserId, err := rdb.Get(fmt.Sprintf("reservation-lock-%s", reserveParkLocationRequest.Id)).Result()

		if err != nil {
			log.Printf("unexpected error %v", err)
			return err, nil
		}

		if reservedUserId != clientKafkaRequestMessage.ClientId {
			return nil, &ErrorMessage{Code: "ui.messages.reservation.location_already_reserved_for_another_user", Message: "location already reserved for other user"}
		}

	*/

	var count int64 = 0
	db.QueryRow("select count(*) from pa_locations where status = $1 and id = $2;", 1, reserveParkLocationRequest.Id).Scan(&count)

	if count > 0 {
		return nil, &ErrorMessage{Code: "ui.messages.reservation.location_already_reserved", Message: "location already reserved"}
	}

	result, err := db.Exec("update pa_locations set status = $1, assigned_user_id = $2 where id = $3;", 1, clientKafkaRequestMessage.ClientId, reserveParkLocationRequest.Id)

	if err != nil {
		log.Printf("db.exec unexpected error %v", err)
		return err, nil
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Printf("rowsAffected - unexpected error %v", err)
		return err, nil
	}

	return err, rowsAffected

}

type ScheduleParkLocationAvailabilityOperation struct {
}

func (o *ScheduleParkLocationAvailabilityOperation) Do(data interface{}) (error, interface{}) {

	clientKafkaRequestMessage := data.(ClientKafkaRequestMessage)

	scheduleParkLocationAvailabiltyRequestData, err := json.Marshal(clientKafkaRequestMessage.Data)

	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	var scheduleParkLocationAvailabiltyRequest ScheduleParkLocationAvailabiltyRequest

	err = json.Unmarshal(scheduleParkLocationAvailabiltyRequestData, &scheduleParkLocationAvailabiltyRequest)
	if err != nil {
		log.Printf("unexpected error %v", err)
		return err, nil
	}

	if scheduleParkLocationAvailabiltyRequest.ScheduleType == ById {
		var location LocationEntity
		row := db.QueryRow("select status,assigned_user_id from pa_locations where id = $1;", scheduleParkLocationAvailabiltyRequest.Id)
		err := row.Err()
		if err != nil {
			log.Printf("unexpected error %v", err)
			return err, nil
		}

		err = row.Scan(&location.Status, &location.AssignedUserId)

		if err != nil {
			log.Printf("unexpected error %v", err)
			return err, nil
		}

		if !location.AssignedUserId.Valid || location.AssignedUserId.String != clientKafkaRequestMessage.ClientId {
			return nil, &ErrorMessage{Code: "ui.messages.reservation.invalid_schedule_location_availability_request", Message: "invalid schedule location availability request"}
		}

		if location.Status == 2 {
			return nil, &ErrorMessage{Code: "ui.messages.reservation.location_already_scheduled", Message: "location already scheduled for availablity"}
		}

		result, err := db.Exec("update pa_locations set status = $1, scheduled_available_time = $2 where id = $3 and assigned_user_id = $4;", 2, scheduleParkLocationAvailabiltyRequest.ScheduledTime, scheduleParkLocationAvailabiltyRequest.Id, clientKafkaRequestMessage.ClientId)

		if err != nil {
			log.Printf("db.exec unexpected error %v", err)
			return err, nil
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			log.Printf("rowsAffected - unexpected error %v", err)
			return err, nil
		}

		return err, rowsAffected
	} else {
		result, err := db.Exec("INSERT INTO pa_locations (id,status,geog,owner_id,scheduled_available_time) VALUES($1,$2, ST_MakePoint($3, $4)::geography,$5)", uuid.New().String(), 2, scheduleParkLocationAvailabiltyRequest.Longitude, scheduleParkLocationAvailabiltyRequest.Latitude, clientKafkaRequestMessage.ClientId, scheduleParkLocationAvailabiltyRequest.ScheduledTime)

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

}

var operations map[string]IOperation = make(map[string]IOperation, 0)

var db *sql.DB
var rdb *redis.Client
var reservationLockDurationSeconds int32 = 15
var reservationLockDurationTime time.Duration = time.Second * 15
var maxSearchResult int32 = 3
var maxDistance float64 = 5000

func main() {

	operations["get_locations_nearby"] = &FindLocationsNearbyOperation{}
	operations["create_park_location"] = &CreateParkLocationOperation{}
	operations["reserve_park_location"] = &ReserveParkLocationOperation{}
	operations["schedule_park_location_availability"] = &ScheduleParkLocationAvailabilityOperation{}

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

	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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
					IsSuccess:     true,
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

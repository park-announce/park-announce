package service

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/global"
	"github.com/segmentio/kafka-go"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

type SocketConnection struct {
	connection *websocket.Conn
}

func (s *SocketService) CreateSocketConnection(ctx *gin.Context, user entity.User, hub *SocketHub) {

	connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("error :", err)
		panic(err)
	}

	client := &SocketClient{conn: &SocketConnection{connection: connection}, user: user}
	hub.register <- client

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9092"),
		Topic:                  "client_request",
		AllowAutoTopicCreation: true,
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Println("failed to close writer:", err)
		}
	}()

	for {
		// Read message from browser
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("error :", err)
			return
		}

		var socketMessage entity.SocketMessage
		deSerializationError := json.Unmarshal(msg, &socketMessage)
		if deSerializationError != nil {
			log.Println(deSerializationError)
			continue
		}

		//check is any transaction exist in redis with the same transaction id for this operation and this client id
		transactionUniqueKey := fmt.Sprintf("transaction|%s|%s|%s", user.Id, socketMessage.Operation, socketMessage.TransactionId)
		isNewTransactionCanBeStarted, err := s.redisClient.SetNX(transactionUniqueKey, 1, time.Second*5)
		if err != nil {
			log.Println("error :", err)
			return
		}

		if !isNewTransactionCanBeStarted {
			log.Printf("isNewTransactionCanBeStarted is false. client id : %s, operation name : %s, transction id : %s ", user.Id, socketMessage.Operation, socketMessage.TransactionId)
			return
		}

		clientKafkaRequestMessage := &entity.ClientKafkaRequestMessage{
			ClientId:        user.Id,
			ApiId:           fmt.Sprintf("%d", global.GetInstanceId()),
			TransactionId:   socketMessage.TransactionId,
			Operation:       socketMessage.Operation,
			Data:            socketMessage.Data,
			Timeout:         socketMessage.Timeout,
			TransactionTime: time.Now().Unix(),
		}

		data, err := json.Marshal(clientKafkaRequestMessage)

		if err != nil {
			log.Println("error :", err)
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
				log.Println("error :", err)
			}
			break
		}

	}
}

func (client *SocketClient) SendMessage(message []byte) {

	err := client.conn.connection.WriteMessage(websocket.TextMessage, message)

	if err != nil {
		log.Println("error :", err)
		panic(err)
	}
}

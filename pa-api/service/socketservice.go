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
		panic(err)
	}

	client := &SocketClient{conn: &SocketConnection{connection: connection}, send: make(chan []byte, 256), user: user}
	hub.register <- client

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
		_, msg, err := connection.ReadMessage()
		if err != nil {
			return
		}

		var socketMessage entity.SocketMessage
		deSerializationError := json.Unmarshal(msg, &socketMessage)
		if deSerializationError != nil {
			log.Println(deSerializationError)
			continue
		}

		clientKafkaRequestMessage := &entity.ClientKafkaRequestMessage{
			ClientId:      socketMessage.ClientId,
			ApiId:         fmt.Sprintf("%d", global.GetInstanceId()),
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
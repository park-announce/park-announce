package service

import (
	"time"

	"github.com/park-announce/pa-api/entity"
)

type SocketClient struct {
	// The websocket connection.
	conn *SocketConnection

	// Buffered channel of outbound messages.
	send chan []byte

	user entity.User
}

type SocketHub struct {

	// Registered clients.
	clients map[string]*SocketClient

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *SocketClient

	// Unregister requests from clients.
	unregister chan *SocketClient
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewSocketHub() *SocketHub {
	return &SocketHub{
		broadcast:  make(chan []byte),
		register:   make(chan *SocketClient),
		unregister: make(chan *SocketClient),
		clients:    make(map[string]*SocketClient),
	}
}

func (h *SocketHub) Broadcast() {

	for {
		select {

		case message := <-h.broadcast:
			for key, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, key)
				}
			}
		}
	}
}

func (h *SocketHub) SendMessage(id string, message string) {

	if client, ok := h.clients[id]; ok {
		client.send <- []byte(message)
	}
}

func (h *SocketHub) SendMessageIfClientExist(id string, message []byte) bool {

	if client, ok := h.clients[id]; ok {
		client.send <- message
		return true
	}
	return false
}

func (h *SocketHub) IsClientExist(id string) bool {

	_, ok := h.clients[id]

	return ok
}

func (h *SocketHub) UnRegister() {

	for {
		select {

		case client := <-h.unregister:
			if _, ok := h.clients[client.user.Id]; ok {
				delete(h.clients, client.user.Id)
				close(client.send)
			}
		}
	}
}

func (h *SocketHub) Register() {

	for {
		select {
		case client := <-h.register:
			h.clients[client.user.Id] = client
		}
	}
}

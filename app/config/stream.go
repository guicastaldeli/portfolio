package config

import (
	"log"
	"main/message"
	"main/server"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	*server.Server
}

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: 10 * time.Second,
	}
	s *server.Server
)

// Broadcast
func (s *Server) SendBroadcast(message message.Message) {
	s.Broadcast <- message
}

func (s *Server) SendBroacastToChannel(channel string, message message.Message) {
	message.Channel = channel
	s.Broadcast <- message
}

// Send to Client
func (s *Server) Send(clientId string, message message.Message) {
	s.Mutex.RLock()
	client, ok := s.Clients[clientId]
	s.Mutex.RUnlock()

	if ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("Channel us full")
		}
	}
}

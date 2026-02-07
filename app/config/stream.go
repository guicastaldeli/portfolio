package config

import (
	"log"
	"main/message"
	"main/ws"
)

type Server struct {
	*ws.Server
}

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

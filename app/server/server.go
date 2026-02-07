package server

import (
	"log"
	"main/message"
	"sync"
)

type Server struct {
	Clients     map[string]*Client
	Broadcast   chan message.Message
	Register    chan *Client
	Unregister  chan *Client
	Subscribe   chan message.Subscription
	Unsubscribe chan message.Subscription
	Mutex       sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Clients:     make(map[string]*Client),
		Broadcast:   make(chan message.Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Subscribe:   make(chan message.Subscription),
		Unsubscribe: make(chan message.Subscription),
	}
}

func (s *Server) Start() {
	log.Println("Websocket server started!")

	for {
		select {
		// Connected
		case client := <-s.Register:
			s.Mutex.Lock()
			s.Clients[client.Id] = client
			s.Mutex.Unlock()

			log.Printf("Client %s connected. Total: %d", client.Id, len(s.Clients))
		// Disconnected
		case client := <-s.Unregister:
			s.Mutex.Lock()
			if c, ok := s.Clients[client.Id]; ok {
				delete(s.Clients, client.Id)
				close(c.Send)
			}
			s.Mutex.Unlock()
			log.Printf("Client %s disconnected; Total: %d", client.Id, len(s.Clients))
		// Subscribe
		case sub := <-s.Subscribe:
			s.Mutex.Lock()
			if client, ok := s.Clients[sub.ClientId]; ok {
				if client.Channels == nil {
					client.Channels = make(map[string]bool)
				}
				client.Channels[sub.Channel] = true
			}
			s.Mutex.Unlock()
		// Unsubscribe
		case unsub := <-s.Unsubscribe:
			s.Mutex.Lock()
			if client, ok := s.Clients[unsub.ClientId]; ok {
				delete(client.Channels, unsub.Channel)
			}
			s.Mutex.Unlock()
		// Broadcast
		case message := <-s.Broadcast:
			s.Mutex.RLock()
			for _, client := range s.Clients {
				if message.Channel != "" {
					if !client.Channels[message.Channel] {
						continue
					}
				}

				select {
				case client.Send <- message:
					log.Printf("Client %s: Message sent", client.Id)
				default:
					close(client.Send)
					delete(s.Clients, client.Id)
				}
			}
			s.Mutex.RUnlock()
		}

	}
}

func Run() *Server {
	server := NewServer()
	go server.Start()

	log.Println("Server starting on :3000")
	return server
}

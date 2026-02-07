package server

import (
	"log"
	"main/config"
	"net/http"
	"sync"
)

type Server struct {
	clients     map[string]*Client
	broadcast   chan config.Message
	register    chan *Client
	unregister  chan *Client
	subscribe   chan config.Subscription
	unsubscribe chan config.Subscription
	mutex       sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		clients:     make(map[string]*Client),
		broadcast:   make(chan config.Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		subscribe:   make(chan config.Subscription),
		unsubscribe: make(chan config.Subscription),
	}
}

func (s *Server) Start() {
	log.Println("Websocket server started!")

	for {
		select {
		// Connected
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client.Id] = client
			s.mutex.Unlock()

			log.Println("Client %s connected. Total: %d", client.Id, len(s.clients))
		// Disconnected
		case client := <-s.unregister:
			s.mutex.Lock()
			if c, ok := s.clients[client.Id]; ok {
				delete(s.clients, client.Id)
				close(c.Send)
			}
			s.mutex.Unlock()
			log.Println("Client %s disconnected; Total: %d", client.Id, len(s.clients))
		// Subscribe
		case sub := <-s.subscribe:
			s.mutex.Lock()
			if client, ok := s.clients[sub.ClientId]; ok {
				if client.Channels == nil {
					client.Channels = make(map[string]bool)
				}
				client.Channels[sub.Channel] = true
			}
			s.mutex.Unlock()
		// Unsubscribe
		case unsub := <-s.unsubscribe:
			s.mutex.Lock()
			if client, ok := s.clients[unsub.ClientId]; ok {
				delete(client.Channels, unsub.Channel)
			}
			s.mutex.Unlock()
		// Broadcast
		case message := <-s.broadcast:
			s.mutex.RLock()
			for _, client := range s.clients {
				if message.Channel != "" {
					if !client.Channels[message.Channel] {
						continue
					}
				}

				select {
				case client.Send <- message:
					log.Println("Client %s: Message sent", client.Id)
				default:
					close(client.Send)
					delete(s.clients, client.Id)
				}
			}
			s.mutex.RUnlock()
		}

	}
}

func Run() {
	config.Setup()
	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

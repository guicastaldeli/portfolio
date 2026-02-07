package server

import (
	"encoding/json"
	"main/message"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id       string
	Conn     *websocket.Conn
	Send     chan message.Message
	Channels map[string]bool
}

// Generate Client Id
func GenerateClientId() string {
	return "client_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(num int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, num)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// Get Connected Clients
func (s *Server) GetConnectedClients() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return len(s.Clients)
}

// Handle Client Message
func (s *Server) HandleClientMessage(client *Client, rawMessage []byte) {
	var msg message.Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		client.Send <- message.Message{
			Type:  "error",
			Error: "Invalid message format",
		}
		return
	}

	switch msg.Type {
	// Subscribe
	case "subscribe":
		if msg.Channel == "" {
			client.Send <- message.Message{
				Type:  "error",
				Error: "Channel name required",
			}
			return
		}

		s.Subscribe <- message.Subscription{
			ClientId: client.Id,
			Channel:  msg.Channel,
		}

		client.Send <- message.Message{
			Type: "subscribed",
			Data: map[string]interface{}{
				"channel": msg.Channel,
			},
		}
	// Unsubscribe
	case "unsubscribe":
		if msg.Channel == "" {
			client.Send <- message.Message{
				Type:  "error",
				Error: "Channel name required",
			}
			return
		}

		s.Unsubscribe <- message.Subscription{
			ClientId: client.Id,
			Channel:  msg.Channel,
		}

		client.Send <- message.Message{
			Type: "unsubscribed",
			Data: map[string]interface{}{
				"channel": msg.Channel,
			},
		}
	// Ping :)))
	case "ping":
		client.Send <- message.Message{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
			},
		}
	default:
		client.Send <- message.Message{
			Type: "echo",
			Data: msg.Data,
		}
	}
}

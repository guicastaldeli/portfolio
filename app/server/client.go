package server

import (
	"encoding/json"
	"main/config"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id       string
	Conn     *websocket.Conn
	Send     chan config.Message
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
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.clients)
}

// Handle Client Message
func (s *Server) handleClientMessage(client *Client, rawMessage []byte) {
	var msg config.Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		client.Send <- config.Message{
			Type:  "error",
			Error: "Invalid message format",
		}
		return
	}

	switch msg.Type {
	// Subscribe
	case "subscribe":
		if msg.Channel == "" {
			client.Send <- config.Message{
				Type:  "error",
				Error: "Channel name required",
			}
			return
		}

		s.subscribe <- config.Subscription{
			ClientId: client.Id,
			Channel:  msg.Channel,
		}

		client.Send <- config.Message{
			Type: "subscribed",
			Data: map[string]interface{}{
				"channel": msg.Channel,
			},
		}
	// Unsubscribe
	case "unsubscribe":
		if msg.Channel == "" {
			client.Send <- config.Message{
				Type:  "error",
				Error: "Channel name required",
			}
			return
		}

		s.unsubscribe <- config.Subscription{
			ClientId: client.Id,
			Channel:  msg.Channel,
		}

		client.Send <- config.Message{
			Type: "unsubscribed",
			Data: map[string]interface{}{
				"channel": msg.Channel,
			},
		}
	// Ping :)))
	case "ping":
		client.Send <- config.Message{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
			},
		}
	default:
		client.Send <- config.Message{
			Type: "echo",
			Data: msg.Data,
		}
	}
}

package config

import (
	"encoding/json"
	"log"
	"main/message"
	"main/server"
	"main/ws"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error!: %v", err)
		return
	}

	clientId := server.GenerateClientId()
	client := &server.Client{
		Id:       clientId,
		Conn:     conn,
		Send:     make(chan message.Message, 256),
		Channels: make(map[string]bool),
	}

	s.Register <- client
	go s.writePump(client)
	go s.readPump(client)

	client.Send <- message.Message{
		Type: "connected",
		Data: map[string]interface{}{
			"clientId":  clientId,
			"timestamp": time.Now().Unix(),
			"message":   "Connected to API",
		},
	}
}

// Write Pump
func (s *Server) writePump(client *server.Client) {
	defer func() {
		client.Conn.Close()
		s.Unregister <- client
	}()

	for message := range client.Send {
		w, err := client.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		if err := json.NewEncoder(w).Encode(message); err != nil {
			log.Printf("Error encoding message: %v", err)
			return
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

// Read Pump
func (s *Server) readPump(client *server.Client) {
	defer func() {
		s.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("Websocket error: %v", err)
			}
			break
		}

		s.HandleClientMessage(client, message)
	}
}

package config

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"main/server"

	"github.com/gorilla/websocket"
)

func (s *server.Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error!: %v", err)
	}

	clientId = GenerateClientId()

	client := server.Client{
		Id: clientId,
		Conn: conn,
		Send: make(chan Message, 256),
		Channels: make(make[string]bool),
	}

	s.register <- client
	go s.writePump(client)
	go s.readPump(client)

	client.Send <- Message{
		Type: "connected",
		Data: map[string]interface{}{
			"clientId": clientId,
			"timestamp": time.Now().Unix()
			"message": "Connected to API"
		},
	}
}
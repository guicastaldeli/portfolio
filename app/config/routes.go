package config

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Endpoint
func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client Connected...")

	reader(ws)
}

// Setup
func Setup() {
	InitIndex()

	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/helloWs", HelloWs)
}

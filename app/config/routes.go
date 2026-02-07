package config

import (
	"log"
	"main/api"
	"main/server"
	"main/ws"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	wsServer *Server
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
	ws.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client Connected...")

	reader(ws)
}

// Setup
func Setup(s *server.Server) {
	wsServer = &Server{s}

	InitIndex()

	http.HandleFunc("/ws", wsServer.HandleWebSocket)
	http.HandleFunc("/time-stream", api.TimeStreamHandler)
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/helloWs", HelloWs)
}

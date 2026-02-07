package ws

import (
	"main/server"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

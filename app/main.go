package main

import (
	"log"
	"main/config"
	"main/db"
	"main/server"
	"main/ws"
	"net/http"
)

func main() {
	cfg := db.Init()
	if err := db.InitDb(cfg); err != nil {
		log.Fatal("Failed to initialize database!", err)
	}
	defer db.CloseDb()

	serverInstance := server.Run()
	wsServer := &ws.Server{
		Server: serverInstance,
	}
	config.Setup(wsServer)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("HTTP server failed to start: ", err)
	}
}

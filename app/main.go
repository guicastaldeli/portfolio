package main

import (
	"log"
	"main/config"
	"main/server"
	"net/http"
)

func main() {
	wsServer := server.Run()
	config.Setup(wsServer)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("HTTP server failed to start: ", err)
	}
}

package main

import (
	"log"
	"main/server"
	"net/http"
)

func main() {
	server.Run()
	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

package main

import (
	"main/config"
	"main/server"
)

func main() {
	config.Setup()
	server.Run()
}

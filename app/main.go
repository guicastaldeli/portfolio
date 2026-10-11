package main

import (
	"log"
	"main/auth"
	"main/config"
	"main/db"
	"main/server"
	"main/ws"
	"net/http"
	"strings"
)

func logs() {
	serverUrl := config.GetEnv("SERVER_INIT_URL")
	apiUrl := config.GetEnv("API_URL")
	webUrl := config.GetEnv("WEB_URL")

	lineStr := strings.Repeat("=", 80)
	emptyStr := strings.Repeat(" ", 20)
	log.Println(lineStr)
	log.Println(" ")
	log.Printf("%s Server starting on %s", emptyStr, serverUrl)
	log.Println(" ")
	log.Println(lineStr)

	log.Printf("Api URL: %s", apiUrl)
	log.Printf("Web URL: %s", webUrl)
}

func main() {
	log.SetFlags(0)

	env := ".env"
	if err := config.InitEnv(env); err != nil {
		log.Fatal("Failed to load env config", err)
	}

	auth.Set(
		config.GetEnv("ADMIN_EMAIL"),
		config.GetEnv("ADMIN_PASSWORD"),
		config.GetEnv("SESSION_COOKIE_NAME"),
	)

	config.GenerateTemp()

	serverAddr := config.GetEnv("SERVER_ADDR")
	logs()

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

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("HTTP server failed to start: ", err)
	}
}

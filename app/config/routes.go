package config

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello!")
}

func helloWs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello WebSocket!")
}

func Setup() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/ws", helloWs)
}

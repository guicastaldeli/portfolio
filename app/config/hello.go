package config

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello!")
}

func HelloWs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello WebSocket!")
}

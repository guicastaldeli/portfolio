package config

import (
	"net/http"
)

func InitIndex() {
	var filePath = "./server/index.html"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	})
}

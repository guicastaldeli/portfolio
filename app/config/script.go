package config

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func InitScripts() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Could not get working directory: %v", err)
		wd = "."
	}

	var jsPath string
	var cssPath string

	// Check if we're in root (with app subdirectory) or in app directory
	appPath := filepath.Join(wd, "app")
	if _, err := os.Stat(appPath); err == nil {
		// We're in root, use app/...
		jsPath = filepath.Join(wd, "app", "server", ".out", "main.js")
		cssPath = filepath.Join(wd, "app", "server", ".styles", "client.css")
	} else {
		// We're probably in app directory or local dev
		jsPath = filepath.Join(wd, "server", ".out", "main.js")
		cssPath = filepath.Join(wd, "server", ".styles", "client.css")
	}

	log.Printf("JS path: %s", jsPath)
	log.Printf("CSS path: %s", cssPath)

	// Serve main.js
	http.HandleFunc("/scripts/main.js", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving main.js")

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.ServeFile(w, r, jsPath)
	})

	// Serve client.css
	http.HandleFunc("/styles/client.css", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving client.css")

		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.ServeFile(w, r, cssPath)
	})

	// If you need to serve other files from .out directory (like server/index.js)
	http.HandleFunc("/scripts/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the file path after /scripts/
		requestedPath := r.URL.Path[len("/scripts/"):]

		var fullPath string
		if _, err := os.Stat(appPath); err == nil {
			fullPath = filepath.Join(wd, "app", "server", ".out", requestedPath)
		} else {
			fullPath = filepath.Join(wd, "server", ".out", requestedPath)
		}

		log.Printf("Serving script: %s from %s", requestedPath, fullPath)

		// Check if file exists
		if _, err := os.Stat(fullPath); err != nil {
			log.Printf("File not found: %s", fullPath)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.ServeFile(w, r, fullPath)
	})
}

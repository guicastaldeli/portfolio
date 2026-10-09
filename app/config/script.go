package config

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func InitScripts() {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Could not get working directory: %v", err)
		wd = "."
	}
	appPath := filepath.Join(wd, "app")
	outDir := filepath.Join(wd, "..", ".out")

	log.Printf("** Working directory: %s", wd)
	log.Printf("** Out directory: %s", outDir)

	// Scripts
	http.HandleFunc("/scripts/", func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path[len("/scripts/"):]
		fullPath := filepath.Join(outDir, requestedPath)
		log.Printf("Serving script: %s from %s", requestedPath, fullPath)

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
	// Styles
	http.HandleFunc("/styles/", func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path[len("/styles/"):]

		var fullPath string
		if _, err := os.Stat(appPath); err == nil {
			fullPath = filepath.Join(wd, "app", "styles", requestedPath)
		} else {
			fullPath = filepath.Join(wd, "styles", requestedPath)
		}

		log.Printf("Serving style: %s from %s", requestedPath, fullPath)

		if _, err := os.Stat(fullPath); err != nil {
			log.Printf("File not found: %s", fullPath)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, fullPath)
	})
}

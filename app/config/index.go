package config

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func InitIndex() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Could not get working directory: %v", err)
		wd = "."
	}

	log.Printf("Working directory: %s", wd)

	// Build paths that work both locally and on Render
	var indexPath string
	var projectEditorPath string
	var staticDir string

	// Check if we're in root (with app subdirectory) or in app directory
	appPath := filepath.Join(wd, "app")
	if _, err := os.Stat(appPath); err == nil {
		// We're in root, use app/server/...
		indexPath = filepath.Join(wd, "app", "server", "index.html")
		projectEditorPath = filepath.Join(wd, "app", "server", "project-editor.html")
		staticDir = filepath.Join(wd, "app", "server")
	} else {
		// We're probably in app directory or local dev
		indexPath = filepath.Join(wd, "server", "index.html")
		projectEditorPath = filepath.Join(wd, "server", "project-editor.html")
		staticDir = filepath.Join(wd, "server")
	}

	log.Printf("Index path: %s", indexPath)
	log.Printf("Project editor path: %s", projectEditorPath)
	log.Printf("Static directory: %s", staticDir)

	// Verify files exist
	if _, err := os.Stat(indexPath); err != nil {
		log.Printf("ERROR: index.html not found at %s", indexPath)
	}
	if _, err := os.Stat(projectEditorPath); err != nil {
		log.Printf("ERROR: project-editor.html not found at %s", projectEditorPath)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Handle editor route
		if r.URL.Path == "/editor" || r.URL.Path == "/editor/" {
			log.Printf("Serving editor page")
			http.ServeFile(w, r, projectEditorPath)
			return
		}

		// Serve root page
		if r.URL.Path == "/" {
			log.Printf("Serving index page")
			http.ServeFile(w, r, indexPath)
			return
		}

		// Serve static files (CSS, JS, etc.)
		// Check if file exists in static directory
		filePath := filepath.Join(staticDir, r.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			log.Printf("Serving static file: %s", filePath)
			http.ServeFile(w, r, filePath)
			return
		}

		// If nothing matches, serve index (SPA fallback)
		log.Printf("Fallback to index page")
		http.ServeFile(w, r, indexPath)
	})
}

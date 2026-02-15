package config

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	var serverDir string

	// Check if we're in root (with app subdirectory) or in app directory
	appPath := filepath.Join(wd, "app")
	if _, err := os.Stat(appPath); err == nil {
		// We're in root, use app/...
		indexPath = filepath.Join(wd, "app", "server", "index.html")
		projectEditorPath = filepath.Join(wd, "app", "server", "project-editor.html")
		serverDir = filepath.Join(wd, "app", "server")
	} else {
		// We're probably in app directory or local dev
		indexPath = filepath.Join(wd, "server", "index.html")
		projectEditorPath = filepath.Join(wd, "server", "project-editor.html")
		serverDir = filepath.Join(wd, "server")
	}

	log.Printf("Index path: %s", indexPath)
	log.Printf("Project editor path: %s", projectEditorPath)
	log.Printf("Server directory: %s", serverDir)

	// Verify files exist
	if _, err := os.Stat(indexPath); err != nil {
		log.Printf("ERROR: index.html not found at %s", indexPath)
	}
	if _, err := os.Stat(projectEditorPath); err != nil {
		log.Printf("ERROR: project-editor.html not found at %s", projectEditorPath)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Handle editor route
		if r.URL.Path == "/editor" || r.URL.Path == "/editor/" {
			log.Printf("Serving editor page")
			http.ServeFile(w, r, projectEditorPath)
			return
		}

		// Serve root page
		if r.URL.Path == "/" {
			log.Printf("Serving index page")
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, indexPath)
			return
		}

		// Handle paths that start with /. (like ./.out/main.js, ./.styles/client.css)
		urlPath := r.URL.Path
		if strings.HasPrefix(urlPath, "/.") {
			// Remove leading /. to get the actual path
			urlPath = strings.TrimPrefix(urlPath, "/")
		}

		// Build full file path
		filePath := filepath.Join(serverDir, urlPath)

		log.Printf("Looking for file: %s", filePath)

		// Check if file exists
		if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() {
			// Set proper MIME types
			if strings.HasSuffix(filePath, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			} else if strings.HasSuffix(filePath, ".css") {
				w.Header().Set("Content-Type", "text/css")
			}

			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			log.Printf("Serving static file: %s", filePath)
			http.ServeFile(w, r, filePath)
			return
		}

		// If file not found, log it
		log.Printf("File not found: %s (looked at %s)", r.URL.Path, filePath)

		// For other paths that don't exist, serve index (SPA fallback)
		log.Printf("Fallback to index page")
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, indexPath)
	})
}

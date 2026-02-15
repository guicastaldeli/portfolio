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

	// Create a file server for the server directory
	fs := http.FileServer(http.Dir(serverDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Handle editor route
		if r.URL.Path == "/editor" || r.URL.Path == "/editor/" {
			log.Printf("Serving editor page")
			w.Header().Set("Content-Type", "text/html")
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

		// Try to serve as static file
		filePath := filepath.Join(serverDir, r.URL.Path)
		log.Printf("Looking for file: %s", filePath)

		// Check if file exists
		if fileInfo, err := os.Stat(filePath); err == nil {
			// If it's a directory, don't serve it
			if fileInfo.IsDir() {
				log.Printf("Path is a directory, serving index fallback")
				w.Header().Set("Content-Type", "text/html")
				http.ServeFile(w, r, indexPath)
				return
			}

			// Set proper MIME types based on file extension
			contentType := getContentType(filePath)
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}

			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			log.Printf("Serving static file: %s with Content-Type: %s", filePath, contentType)

			// Use the file server to handle the request
			fs.ServeHTTP(w, r)
			return
		}

		// If file not found, log it and serve index (SPA fallback)
		log.Printf("File not found: %s (looked at %s)", r.URL.Path, filePath)
		log.Printf("Fallback to index page")
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, indexPath)
	})
}

// Helper function to determine content type based on file extension
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".js", ".mjs":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".html", ".htm":
		return "text/html"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".otf":
		return "font/otf"
	default:
		return ""
	}
}

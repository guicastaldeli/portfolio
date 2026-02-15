package config

import (
	"net/http"
	"os"
	"path/filepath"
)

func InitIndex() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		// Fallback to relative paths if we can't get working directory
		wd = "."
	}

	// Build paths that work both locally and on Render
	var indexPath string
	var projectEditorPath string

	// Check if we're in root (with app subdirectory) or in app directory
	if _, err := os.Stat(filepath.Join(wd, "app")); err == nil {
		// We're in root, use app/server/...
		indexPath = filepath.Join(wd, "app", "server", "index.html")
		projectEditorPath = filepath.Join(wd, "app", "server", "project-editor.html")
	} else {
		// We're probably in app directory or local dev
		indexPath = filepath.Join(wd, "server", "index.html")
		projectEditorPath = filepath.Join(wd, "server", "project-editor.html")
	}

	fs := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "." + r.URL.Path

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Handle editor route
		if r.URL.Path == "/editor" || r.URL.Path == "/editor/" {
			http.ServeFile(w, r, projectEditorPath)
			return
		}

		// Serve static files if they exist
		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		// For root path or any other path, serve the server page
		http.ServeFile(w, r, indexPath)
	})
}

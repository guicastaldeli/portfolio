package config

import (
	"net/http"
	"os"
)

func InitIndex() {
	var indexPath = "./server/index.html"
	var projectEditorPath = "./server/project-editor.html"
	fs := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "." + r.URL.Path

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if r.URL.Path == "/editor" || r.URL.Path == "/editor/" {
			http.ServeFile(w, r, projectEditorPath)
			return
		}
		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, indexPath)
	})
}

package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"main/db"
	"main/message"
	"main/ws"
	"net/http"
	"strconv"
	"strings"
)

// Get Projects
func GetAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectDb, err := db.GetDb("project")
	if err != nil {
		log.Printf("Project database error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := projectDb.Query(db.Q(db.GetAllProjects))
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	projects := []message.Project{}
	for rows.Next() {
		var p message.Project
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Desc,
			&p.Repo,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning project: %v", err)
			continue
		}

		p.Media = getProjectMedia(p.Id)
		p.Links = getProjectLinks(p.Id)
		projects = append(projects, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if idStr == "" || idStr == "/" {
		http.Error(w, "Invalid project Id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project Id", http.StatusBadRequest)
		return
	}

	projectDb, err := db.GetDb("project")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var p message.Project
	row := projectDb.QueryRow(db.Q(db.GetProjectById), id)

	err = row.Scan(
		&p.Id,
		&p.Name,
		&p.Desc,
		&p.Repo,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.Media = getProjectMedia(p.Id)
	p.Links = getProjectLinks(p.Id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Create Project
func CreateProjectHandler(wsServer *ws.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req message.CreateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		database, err := db.GetDb("project")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx, err := database.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec(
			db.Q(db.InsertProject),
			req.Name,
			req.Desc,
			req.Repo,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		projectId, _ := res.LastInsertId()
		for _, photo := range req.Photos {
			_, err := tx.Exec(
				db.Q(db.InsertMedia),
				projectId,
				"photo",
				photo,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		for _, video := range req.Videos {
			_, err := tx.Exec(
				db.Q(db.InsertMedia),
				projectId,
				"video",
				video,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		for _, link := range req.Links {
			_, err := tx.Exec(
				db.Q(db.InsertLink),
				projectId,
				link.Name,
				link.URL,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wsServer.Broadcast <- message.Message{
			Type:    "project_created",
			Channel: "projects",
			Data: map[string]interface{}{
				"id":   projectId,
				"name": req.Name,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      projectId,
			"message": "Project created successfully",
		})
	}
}

// Update Project
func UpdateProjectHandler(wsServer *ws.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			log.Printf("Wrong method: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/")
		log.Printf("Path: %s, ID string: %s", r.URL.Path, idStr)

		if idStr == "" || idStr == "/" {
			log.Printf("Empty ID string")
			http.Error(w, "Invalid project Id", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Invalid ID conversion: %s, error: %v", idStr, err)
			http.Error(w, "Invalid project Id", http.StatusBadRequest)
			return
		}

		log.Printf("Project ID to update: %d", id)

		bodyBytes, _ := io.ReadAll(r.Body)
		log.Printf("Raw request body: %s", string(bodyBytes))

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req message.UpdateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("JSON decode error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		projectDb, err := db.GetDb("project")
		if err != nil {
			log.Printf("Project database connection error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Getting media database connection")
		mediaDb, err := db.GetDb("media")
		if err != nil {
			log.Printf("Media database connection error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Getting links database connection")
		linksDb, err := db.GetDb("links")
		if err != nil {
			log.Printf("Links database connection error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		txProject, err := projectDb.Begin()
		if err != nil {
			log.Printf("Project transaction begin error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer txProject.Rollback()

		txMedia, err := mediaDb.Begin()
		if err != nil {
			log.Printf("Media transaction begin error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer txMedia.Rollback()

		txLinks, err := linksDb.Begin()
		if err != nil {
			log.Printf("Links transaction begin error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer txLinks.Rollback()

		result, err := txProject.Exec(
			db.Q(db.UpdateProject),
			req.Name,
			req.Desc,
			req.Repo,
			id,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		log.Printf("Project Update successful, rows affected: %d", rowsAffected)

		_, err = txMedia.Exec(db.Q(db.DeleteProjectMedia), id)
		if err != nil {
			log.Printf("Delete media error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = txLinks.Exec(db.Q(db.DeleteProjectLinks), id)
		if err != nil {
			log.Printf("Delete links error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for i, photo := range req.Photos {
			log.Printf("Photo %d: %s", i+1, photo)
			_, err := txMedia.Exec(
				db.Q(db.InsertMedia),
				id,
				"photo",
				photo,
			)
			if err != nil {
				log.Printf("Insert photo error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		for i, video := range req.Videos {
			log.Printf("Video %d: %s", i+1, video)
			_, err := txMedia.Exec(
				db.Q(db.InsertMedia),
				id,
				"video",
				video,
			)
			if err != nil {
				log.Printf("Insert video error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		for i, link := range req.Links {
			log.Printf("Link %d: Name='%s', URL='%s'", i+1, link.Name, link.URL)
			_, err := txLinks.Exec(
				db.Q(db.InsertLink),
				id,
				link.Name,
				link.URL,
			)
			if err != nil {
				log.Printf("Insert link error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := txProject.Commit(); err != nil {
			log.Printf("Project transaction commit error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := txMedia.Commit(); err != nil {
			log.Printf("Media transaction commit error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := txLinks.Commit(); err != nil {
			log.Printf("Links transaction commit error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wsServer.Broadcast <- message.Message{
			Type:    "project_updated",
			Channel: "projects",
			Data: map[string]interface{}{
				"id": id,
			},
		}

		log.Printf("WebSocket broadcast sent")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Project updated successfully",
		})
	}
}

// Delete Project
func DeleteProjectHandler(wsServer *ws.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Project Id", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("project", db.Q(db.DeleteProject), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wsServer.Broadcast <- message.Message{
			Type:    "project_deleted",
			Channel: "projects",
			Data: map[string]interface{}{
				"id": id,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Project deleted successfully",
		})
	}
}

func getProjectMedia(projectId int) []message.Media {
	mediaDb, err := db.GetDb("media")
	if err != nil {
		log.Printf("Media database error: %v", err)
		return []message.Media{}
	}

	rows, err := mediaDb.Query(db.Q(db.GetProjectMedia), projectId)
	if err != nil {
		log.Printf("Error loading media: %v", err)
		return []message.Media{}
	}
	defer rows.Close()

	media := []message.Media{}
	for rows.Next() {
		var m message.Media
		err := rows.Scan(
			&m.Id,
			&m.ProjectId,
			&m.Type,
			&m.URL,
		)
		if err != nil {
			log.Printf("Error scanning media: %v", err)
			continue
		}
		media = append(media, m)
	}
	return media
}

func getProjectLinks(projectId int) []message.Link {
	linksDb, err := db.GetDb("links")
	if err != nil {
		log.Printf("Links database error: %v", err)
		return []message.Link{}
	}

	rows, err := linksDb.Query(db.Q(db.GetProjectLinks), projectId)
	if err != nil {
		log.Printf("Error loading links: %v", err)
		return []message.Link{}
	}
	defer rows.Close()

	links := []message.Link{}
	for rows.Next() {
		var l message.Link
		err := rows.Scan(
			&l.Id,
			&l.ProjectId,
			&l.Name,
			&l.URL,
		)
		if err != nil {
			log.Printf("Error scanning link: %v", err)
			continue
		}
		links = append(links, l)
	}
	return links
}

// Handlers
func HandleProjects(wsServer *ws.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllProjectsHandler(w, r)
		case http.MethodPost:
			CreateProjectHandler(wsServer)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func HandleProjectById(wsServer *ws.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetProjectHandler(w, r)
		case http.MethodPut:
			UpdateProjectHandler(wsServer)(w, r)
		case http.MethodDelete:
			DeleteProjectHandler(wsServer)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

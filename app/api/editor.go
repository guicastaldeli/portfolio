package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/db"
	"main/message"
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

	rows, err := db.Query("project", db.Q(db.GetAllProjects))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []message.Project
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
			log.Println("Error scanning project", err)
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

	idStr := strings.TrimPrefix(r.URL.Path, "/api/projects")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project Id", http.StatusBadRequest)
		return
	}

	var p message.Project
	row, err := db.QueryRow(
		"project",
		db.Q(db.GetProjectById),
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = row.Scan(
		&p.Id,
		&p.Desc,
		&p.Repo,
		&p.Name,
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
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	tx, err := database.Begin()
	if err != nil {
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
			db.Q(db.InsertMedia),
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      projectId,
		"message": "Project created successfully",
	})
}

// Update Project
func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/projects")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project Id", http.StatusBadRequest)
		return
	}

	var req message.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database, err := db.GetDb("project")
	if err != nil {
		return
	}

	tx, err := database.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(
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

	_, err = tx.Exec(db.Q(db.DeleteProjectMedia), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec(db.Q(db.DeleteProjectLinks), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, photo := range req.Photos {
		_, err := tx.Exec(
			db.Q(db.InsertMedia),
			id,
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
			id,
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
			db.Q(db.InsertMedia),
			id,
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Project updated successfully",
	})
}

// Delete Project
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Project deleted successfully",
	})
}

func getProjectMedia(projectId int) []message.Media {
	rows, err := db.Query("project", db.Q(db.GetProjectMedia), projectId)
	if err != nil {
		log.Println("Error loading media", err)
		return []message.Media{}
	}
	defer rows.Close()

	var media []message.Media
	for rows.Next() {
		var m message.Media
		rows.Scan(
			&m.Id,
			&m.ProjectId,
			&m.Type,
			&m.URL,
		)
		media = append(media, m)
	}
	return media
}

func getProjectLinks(projectId int) []message.Link {
	rows, err := db.Query("project", db.Q(db.GetProjectLinks), projectId)
	if err != nil {
		log.Println("Error loading links", err)
		return []message.Link{}
	}
	defer rows.Close()

	var links []message.Link
	for rows.Next() {
		var l message.Link
		rows.Scan(
			&l.Id,
			&l.ProjectId,
			&l.Name,
			&l.URL,
		)
		links = append(links, l)
	}
	return links
}

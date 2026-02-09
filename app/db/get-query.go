package db

// Key
type QueryKey string

const (
	// Projects
	GetAllProjects QueryKey = "GET_ALL_PROJECTS"
	GetProjectById QueryKey = "GET_PROJECT_BY_ID"
	InsertProject  QueryKey = "INSERT_PROJECT"
	UpdateProject  QueryKey = "UPDATE_PROJECT"
	DeleteProject  QueryKey = "DELETE_PROJECT"

	// Media
	GetProjectMedia    QueryKey = "GET_PROJECT_MEDIA"
	InsertMedia        QueryKey = "INSERT_MEDIA"
	DeleteProjectMedia QueryKey = "DELETE_PROJECT_MEDIA"

	// Links
	GetProjectLinks    QueryKey = "GET_PROJECT_LINKS"
	InsertLink         QueryKey = "INSERT_LINK"
	DeleteProjectLinks QueryKey = "DELETE_PROJECT_LINKS"
)

// Registry
var QueryRegistry = map[QueryKey]string{
	// Projects
	GetAllProjects: `
		SELECT id, name, description, repo, createdAt, updatedAt
		FROM project
		ORDER BY updatedAt DESC
	`,
	GetProjectById: `
		SELECT id, name, description, repo, createdAt, updatedAt
		FROM project
		WHERE id = ?
	`,
	InsertProject: `
		INSERT INTO project(name, description, repo)
		VALUES(?, ?, ?)
	`,
	UpdateProject: `
		UPDATE project
		SET 
			name = ?,
			description = ?,
			repo = ?,
			updatedAt = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
	DeleteProject: `
		DELETE FROM project WHERE id = ?
	`,

	// Media
	GetProjectMedia: `
		SELECT id, projectId, type, url
		FROM media
		WHERE projectId = ?
	`,
	InsertMedia: `
		INSERT INTO media (projectId, type, url)
		VALUES (?, ?, ?)
	`,
	DeleteProjectMedia: `
		DELETE FROM media WHERE projectId = ?
	`,

	// Links
	GetProjectLinks: `
		SELECT id, projectId, name, url
		FROM links
		WHERE projectId = ?
	`,
	InsertLink: `
		INSERT INTO links(projectId, name, url)
		VALUES (?, ?, ?)
	`,
	DeleteProjectLinks: `
		DELETE FROM links WHERE projectId = ?
	`,
}

// Get Query
func GetQuery(key QueryKey) string {
	query, exists := QueryRegistry[key]
	if !exists {
		panic("Query not found!: " + string(key))
	}
	return query
}

func Q(key QueryKey) string {
	return GetQuery(key)
}

package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Get SQL Files
func getSqlFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(
		path string,
		d fs.DirEntry,
		err error,
	) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".sql") {
			files = append(files, path)
		}

		return nil
	})

	sort.Strings(files)
	return files, err
}

// Execute SQL File
func execSqlFile(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("Failed to execute SQL file %s: %w", path, err)
	}

	return nil
}

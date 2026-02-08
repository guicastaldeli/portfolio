package db

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load Sql Files
func LoadSqlFiles(dir string) error {
	files, err := getSqlFiles(dir)
	if err != nil {
		return fmt.Errorf("failed to read SQL dir: %w", err)
	}
	if len(files) == 0 {
		log.Println("No sql files found!", dir)
		return nil
	}

	sort.Strings(files)
	log.Println("Loading %d SQL files...", len(files))

	for _, file := range files {
		if err := execSqlFile(file); err != nil {
			return fmt.Errorf("failed to execute %s: %w", file, err)
		}
		log.Println("Executed: %s", filepath.Base(file))
	}

	log.Println("All SQL files loaded!")
	return nil
}

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

	return files, err
}

// Execute SQL File
func execSqlFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	stmts := splitStmts(string(content))
	for i, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		_, err := DB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("statement %d failed: %w\n SQL: %s", i+1, err, stmt)
		}
	}

	return nil
}

func splitStmts(content string) []string {
	var stmts []string
	var current strings.Builder
	inString := false
	var stringChar rune

	for i, char := range content {
		switch char {
		case '\'', '"':
			if !inString {
				inString = true
				stringChar = char
			} else if char == stringChar {
				if i > 0 && content[i-1] != '\\' {
					inString = false
				}
			}
			current.WriteRune(char)
		case ';':
			if !inString {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" {
					stmts = append(stmts, stmt)
				}
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		stmts = append(stmts, stmt)
	}

	return stmts
}

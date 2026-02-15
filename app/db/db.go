package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB = make(map[string]*sql.DB)

type Config struct {
	DataDir string
	SrcDir  string
}

// Init
func InitDb(config Config) error {
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}

	sqlFiles, err := getSqlFiles(config.SrcDir)
	if err != nil {
		return fmt.Errorf("failed to get SQL files: %w", err)
	}
	if len(sqlFiles) == 0 {
		log.Println("No SQL files foudn in", config.SrcDir)
	}

	log.Printf("Found %d SQL files, creating dbs...", len(sqlFiles))

	for _, sqlFile := range sqlFiles {
		if err := createDbFromSqlFile(sqlFile, config.DataDir); err != nil {
			return fmt.Errorf("failed to create DB from %s: %w", sqlFile, err)
		}
	}

	log.Println("All databases initialized!")
	return nil
}

func createDbFromSqlFile(sqlFilePath, dataDir string) error {
	baseName := filepath.Base(sqlFilePath)
	dbName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	dbPath := filepath.Join(dataDir, dbName+".db")

	log.Printf("Creating databse: %s from %s", dbPath, baseName)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database %s: %w", dbPath, err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return fmt.Errorf("failed to enable foreign keys for %s: %w", dbPath, err)
	}

	if err := execSqlFile(db, sqlFilePath); err != nil {
		db.Close()
		return fmt.Errorf("failed to execute SQL file %s: %w", sqlFilePath, err)
	}

	DB[dbName] = db
	log.Printf("Databse created: %s", dbPath)
	return nil
}

// Get Db
func GetDb(name string) (*sql.DB, error) {
	db, exists := DB[name]
	if !exists {
		return nil,
			fmt.Errorf("Database '%s' not found in registry", name)
	}
	return db, nil
}

// Close Db
func CloseDb() {
	for name, db := range DB {
		if db != nil {
			db.Close()
			log.Printf("Closed databse: %s", name)
		}
	}
	log.Println("All database connections closed!")
}

func Query(
	dbName string,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	db, err := GetDb(dbName)
	if err != nil {
		return nil, err
	}
	return db.Query(query, args...)
}

func QueryRow(
	dbName string,
	query string,
	args ...interface{},
) (*sql.Row, error) {
	db, err := GetDb(dbName)
	if err != nil {
		return nil, err
	}
	return db.QueryRow(query, args...), nil
}

func Exec(
	dbName string,
	query string,
	args ...interface{},
) (sql.Result, error) {
	db, err := GetDb(dbName)
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

// Init
func Init() Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	appDbPath := filepath.Join(wd, "app", "db")
	if _, err := os.Stat(appDbPath); err == nil {
		return Config{
			DataDir: filepath.Join(wd, "app", "db", "data"),
			SrcDir:  filepath.Join(wd, "app", "db", "src"),
		}
	}

	return Config{
		DataDir: filepath.Join(wd, "db", "data"),
		SrcDir:  filepath.Join(wd, "db", "src"),
	}
}

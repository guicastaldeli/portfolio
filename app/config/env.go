package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Env struct {
	values map[string]string
	mutex  sync.RWMutex
}

var (
	instance *Env
	once     sync.Once
)

// Init
func InitEnv(envDir string) error {
	var err error
	once.Do(func() {
		instance = &Env{
			values: make(map[string]string),
		}
		err = instance.loadEnvFiles(envDir)
	})
	return err
}

// Get
func GetEnv(key string) string {
	if instance == nil {
		log.Fatal("env config not init")
	}

	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	return instance.values[key]
}

func MustGet(key string) string {
	value := GetEnv(key)
	if value == "" {
		log.Fatalf("Required env variable %s not found", key)
	}
	return value
}

// Load
func (c *Env) loadEnvFiles(dir string) error {
	pattern := filepath.Join(dir, ".env*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("Failed to read env dir: %w", err)
	}

	var files []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			log.Printf("Could not stat %s: %v", match, err)
			continue
		}
		if !info.IsDir() {
			files = append(files, match)
		}
	}

	if len(files) == 0 {
		log.Printf("No .env files found in %s - using system environment variables", dir)
		return nil
	}

	log.Printf("Found %d .env file(s)", len(files))

	for _, file := range files {
		if err := c.loadEnvFile(file); err != nil {
			log.Printf("Failed to load %s: %v", file, err)
			continue
		}
		log.Printf("Loaded env file: %s", filepath.Base(file))
	}

	return nil
}

func (c *Env) loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Invalid line %d in %s: %s", lineNum, filename, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = strings.Trim(value, `"'`)

		c.mutex.Lock()
		c.values[key] = value
		c.mutex.Unlock()
	}

	return scanner.Err()
}

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

		instance.loadSystemEnv()
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
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	log.Printf("APP_ENV: %s", appEnv)

	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", appEnv),
	}

	if appEnv == "prod" {
		envFiles = append(envFiles, ".env.production")
	}

	var loadedCount int

	for _, envFile := range envFiles {
		filePath := filepath.Join(dir, envFile)
		if _, err := os.Stat(filePath); err == nil {
			if err := c.loadEnvFile(filePath); err != nil {
				log.Printf("Warning: Failed to load %s: %v", envFile, err)
			} else {
				loadedCount++
				log.Printf("Loaded env file: %s", envFile)
			}
		} else {
			log.Printf("Env file not found: %s (skipping)", envFile)
		}
	}

	if loadedCount == 0 {
		log.Printf("No .env files loaded - using only system environment variables")
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

func (c *Env) loadSystemEnv() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			c.mutex.Lock()
			c.values[parts[0]] = parts[1]
			c.mutex.Unlock()
		}
	}
	log.Printf("Loaded %d system environment variables", len(c.values))
}

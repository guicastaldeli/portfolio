package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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

type Roles string

var Role = struct {
	PRIVATE Roles
	PUBLIC  Roles
}{
	PRIVATE: "private",
	PUBLIC:  "public",
}

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

// Read Role
func readRole(path string) Roles {
	f, err := os.Open(path)
	if err != nil {
		return "public"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() && lines < 20 {
		lines++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			break
		}
		if !strings.Contains(line, "@role:") {
			continue
		}

		parts := strings.SplitN(line, "@role:", 2)
		if len(parts) != 2 {
			continue
		}
		role := strings.TrimSpace(parts[1])
		role = strings.ToLower(role)

		if role == string(Role.PRIVATE) {
			return Role.PRIVATE
		}
		return Role.PUBLIC
	}

	return Role.PUBLIC
}

// Load
func (c *Env) loadEnvFiles(dir string) error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	log.Printf("APP_ENV: %s", appEnv)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read env dir %s: %w", dir, err)
	}

	en := ".env"
	appEnvFile := en + "." + appEnv

	var public, private []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()

		if !strings.HasPrefix(name, en) {
			continue
		}
		if name == KEY_NAME || name == TEMP_NAME {
			continue
		}

		fullPath := filepath.Join(dir, name)
		role := readRole(fullPath)

		if role == Role.PUBLIC {
			if name != en && name != appEnvFile {
				log.Printf("Skipping env file: %s (APP_ENV=%s)", name, appEnv)
				continue
			}

			public = append(public, name)
			continue
		}

		private = append(private, name)
	}

	sort.Strings(public)
	sort.Strings(private)
	ordered := append(public, private...)

	var loaded int
	for _, name := range ordered {
		fullPath := filepath.Join(dir, name)
		if err := c.loadEnvFile(fullPath); err != nil {
			log.Printf("Warning: Failed to load %s: %v", name, err)
		} else {
			loaded++
			log.Printf("Loaded env file: %s", name)
		}
	}
	if loaded == 0 {
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

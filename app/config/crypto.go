package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const KEY_NAME = "env.key"
const TEMP_NAME = "env.temp"

var ekey string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("[env] getwd: %v", err)
		return
	}

	keyPath := filepath.Join(wd, ".env", KEY_NAME)

	b, err := os.ReadFile(keyPath)
	if err != nil {
		log.Printf("[env] cannot read key file %s: %v", keyPath, err)
		return
	}

	ekey = strings.TrimSpace(string(b))
	if ekey == "" {
		log.Printf("[env] key file %s is empty", keyPath)
		return
	}

	//log.Printf("[env] loaded key from %s (%d chars)", keyPath, len(ekey))
}

// Encrypt
func encrypt(plaintext []byte) (string, error) {
	key := sha256.Sum256([]byte(ekey))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Generate Temp
func GenerateTemp() {
	if ekey == "" {
		log.Printf("[env] skipping env.temp — no key loaded")
		return
	}

	appEnv := GetEnv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("[env] getwd: %v", err)
		return
	}

	srcPath := filepath.Join(wd, ".env", ".env."+appEnv)
	plaintext, err := os.ReadFile(srcPath)
	if err != nil {
		log.Printf("[env] cannot read: %s: %v", srcPath, err)
		return
	}

	encrypted, err := encrypt(plaintext)
	if err != nil {
		log.Printf("[env] encrypt: %v", err)
		return
	}

	outPath := filepath.Join(wd, ".env", "env.temp")
	if err := os.WriteFile(outPath, []byte(encrypted), 0644); err != nil {
		log.Printf("[env] write %s: %v", outPath, err)
		return
	}

	log.Printf("[env] wrote encrypted env.temp for env=%s (%d bytes)", appEnv, len(encrypted))
}

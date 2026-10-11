package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	SESSION_COOKIE_NAME string
	ADMIN_EMAIL         string
	ADMIN_PASSWORD      string
)

func Set(sessionCookieName, email, password string) {
	if sessionCookieName != "" {
		SESSION_COOKIE_NAME = sessionCookieName
	}
	if email == "" || password == "" {
		log.Fatal("auth: email or password missing")
	}

	email = ADMIN_EMAIL
	password = ADMIN_PASSWORD
}

// ----------
// Session
// ----------
const sessionTTL = 24 * time.Hour

type sessionEntry struct {
	expiry time.Time
}

var (
	sessions   = map[string]sessionEntry{}
	sessionsMu sync.Mutex
)

func newSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Printf("auth: rand: %v", err)
		return ""
	}
	return hex.EncodeToString(b)
}

func storeSession(token string) {
	sessionsMu.Lock()
	sessions[token] = sessionEntry{expiry: time.Now().Add(sessionTTL)}
	sessionsMu.Unlock()
}

func validSession(token string) bool {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	e, ok := sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(e.expiry) {
		delete(sessions, token)
		return false
	}
	return true
}

func dropSession(token string) {
	sessionsMu.Lock()
	delete(sessions, token)
	sessionsMu.Unlock()
}

// ----------
// Login
// ----------
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Email != ADMIN_EMAIL {
		http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(ADMIN_PASSWORD), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
		return
	}

	token := newSessionToken()
	if token == "" {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	storeSession(token)

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})

	w.WriteHeader(http.StatusNoContent)
}

// ----------
// Logout
// ----------
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if c, err := r.Cookie(SESSION_COOKIE_NAME); err == nil {
		dropSession(c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

// ----------
// Auth
// ----------
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SESSION_COOKIE_NAME)
		if err != nil || !validSession(c.Value) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil || !validSession(c.Value) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Is Authed
func IsAuthed(r *http.Request) bool {
	c, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return false
	}
	return validSession(c.Value)
}

package api

import (
	"encoding/json"
	"log"
	"main/ws"
	"net/http"
	"time"
)

type TimeUpdate struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Formatted string `json:"formatted"`
	Unix      int64  `json:"unix"`
}

func TimeStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		setTimeStream(w, r)
		return
	}

	displayTimeStream(w)
}

// Set
func setTimeStream(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				close(done)
				return
			}
		}
	}()

	for {
		select {
		case t := <-ticker.C:
			timeUpdate := TimeUpdate{
				Type:      "timeUpdate",
				Timestamp: t.Format(time.RFC3339),
				Formatted: t.Format("2006-01-02 15:04:05"),
				Unix:      t.Unix(),
			}

			err := conn.WriteJSON(timeUpdate)
			if err != nil {
				log.Printf("Error %v", err)
				return
			}
		case <-done:
			return
		}
	}
}

// Display
func displayTimeStream(w http.ResponseWriter) {
	now := time.Now()

	data := map[string]interface{}{
		"type":      "currentTime",
		"timestamp": now.Format(time.RFC3339),
		"formatted": now.Format("2006-01-02 15:04:05"),
		"unix":      now.Unix(),
		"timezone":  now.Location().String(),
		"day":       now.Format("Monday"),
		"date":      now.Format("January 2, 2006"),
		"time":      now.Format("3:04:05 PM"),
		"components": map[string]interface{}{
			"year":   now.Year(),
			"month":  int(now.Month()),
			"day":    now.Day(),
			"hour":   now.Hour(),
			"minute": now.Minute(),
			"second": now.Second(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

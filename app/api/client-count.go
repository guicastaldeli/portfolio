package api

import (
	"encoding/json"
	"log"
	"main/ws"
	"net/http"
	"time"
)

type ClientsUpdate struct {
	Type      string   `json:"type"`
	Count     int      `json:"count"`
	Timestamp string   `json:"timestamp"`
	Clients   []string `json:"clients"`
}

func ClientsConnectedHandler(s *ws.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			setClientCount(s, w, r)
			return
		}

		displayClientCount(s, w)
	}
}

// Set
func setClientCount(
	s *ws.Server,
	w http.ResponseWriter,
	r *http.Request,
) {
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
		case <-ticker.C:
			s.Mutex.RLock()
			clientIds := make([]string, 0, len(s.Clients))
			for id := range s.Clients {
				clientIds = append(clientIds, id)
			}
			count := len(s.Clients)
			s.Mutex.RUnlock()

			upgrade := ClientsUpdate{
				Type:      "clientsUpdate",
				Count:     count,
				Timestamp: time.Now().Format(time.RFC3339),
				Clients:   clientIds,
			}

			err := conn.WriteJSON(upgrade)
			if err != nil {
				log.Printf("Error writing to websocket: %v", err)
				return
			}
		case <-done:
			return
		}
	}
}

// Display
func displayClientCount(s *ws.Server, w http.ResponseWriter) {
	s.Mutex.RLock()
	clientIds := make([]string, 0, len(s.Clients))
	clientDetails := make([]map[string]interface{}, 0, len(s.Clients))

	for id, client := range s.Clients {
		clientIds = append(clientIds, id)

		channels := make([]string, 0, len(client.Channels))
		for ch := range client.Channels {
			channels = append(channels, ch)
		}

		clientDetails = append(clientDetails, map[string]interface{}{
			"id":       id,
			"channels": channels,
		})
	}
	count := len(s.Clients)
	s.Mutex.RUnlock()

	data := map[string]interface{}{
		"type":      "clientsSnapshot",
		"count":     count,
		"timestamp": time.Now().Format(time.RFC3339),
		"clients":   clientIds,
		"details":   clientDetails,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

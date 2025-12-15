package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// We keep a store instance here for simplicity
var defaultStore = NewStore()

// CreateUserHandler creates a user (very minimal)
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	userID := uuid.New().String()
	user := User{
		ID:        userID,
		Username:  req.Username,
		CreatedAt: time.Now(),
	}

	defaultStore.CreateUser(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         userID,
		"username":   req.Username,
		"created_at": user.CreatedAt,
	})
}

// WebSocketHandler handles websocket connections
func WebSocketHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Get user ID and room from query parameters
		userID := r.URL.Query().Get("user_id")
		room := r.URL.Query().Get("room")

		if userID == "" {
			userID = uuid.New().String()
		}

		client := &Client{
			Hub:  hub,
			Conn: conn,
			Send: make(chan Message, 256),
			ID:   userID,
			Room: room,
		}

		client.Hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
		go client.ReadPump()
	}
}

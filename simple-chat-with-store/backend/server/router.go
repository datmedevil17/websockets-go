package server

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all the routes for the application
func RegisterRoutes(r *mux.Router, hub *Hub) {
	// API routes
	r.HandleFunc("/api/users", CreateUserHandler).Methods("POST")

	// WebSocket route
	r.HandleFunc("/ws", WebSocketHandler(hub)).Methods("GET")

	// Serve static files (if needed)
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
}

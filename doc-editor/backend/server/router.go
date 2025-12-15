package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, hub *Hub) {
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		id := uuid.New().String()
		hub.store.CreateDocument(id)

		json.NewEncoder(w).Encode(map[string]string{
			"docId": id,
			"link":  "/ws?doc=" + id,
		})
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})
}

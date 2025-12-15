package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, hub *Hub) {
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		hub.store.CreateDocument(id)

		json.NewEncoder(w).Encode(map[string]string{
			"docId": id,
			"link":  "/ws?doc=" + id,
		})
	}).Methods("POST")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})
}

package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {
	docID := r.URL.Query().Get("doc")
	if docID == "" {
		http.Error(w, "doc required", 400)
		return
	}

	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{
		Hub:   h,
		Conn: conn,
		Send: make(chan Message, 256),
		DocID: docID,
	}

	h.register <- client

	go client.WritePump()
	client.ReadPump()
}

package server

import (
	"net/http"

	"github.com/google/uuid"
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientID := uuid.New().String()

	client := &Client{
		Hub:   h,
		Conn:  conn,
		Send:  make(chan Message, 256),
		DocID: docID,
		ID:    clientID,
	}

	h.register <- client

	go client.WritePump()

	// Send identity
	client.Send <- Message{
		Type:     Identity,
		ClientID: clientID,
		DocID:    docID,
	}

	// Send current document state
	doc, ok := h.store.GetDocument(docID)
	if ok {
		client.Send <- Message{
			Type:   SyncDocument,
			DocID:  docID,
			Blocks: doc.Blocks,
		}
	}

	client.ReadPump()
}

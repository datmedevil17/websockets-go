package server

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients map[*Client]bool
	rooms   map[string]map[*Client]bool // roomID -> clients

	register   chan *Client
	unregister chan *Client
	broadcast  chan *MessageEnvelope

	mu sync.RWMutex
}

type MessageEnvelope struct {
	Msg      *Message
	ToRoom   string
	ToClient string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *MessageEnvelope, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.Room != "" {
				if _, ok := h.rooms[client.Room]; !ok {
					h.rooms[client.Room] = make(map[*Client]bool)
				}
				h.rooms[client.Room][client] = true
			}
			h.mu.Unlock()
			log.Printf("client registered: %s room=%s", client.ID, client.Room)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				if client.Room != "" {
					delete(h.rooms[client.Room], client)
				}
			}
			h.mu.Unlock()
			log.Printf("client unregistered: %s", client.ID)

		case envelope := <-h.broadcast:
			if envelope.ToRoom != "" {
				h.mu.RLock()
				clients := h.rooms[envelope.ToRoom]
				for c := range clients {
					// non-blocking send
					select {
					case c.Send <- *envelope.Msg:
					default:
					}
				}
				h.mu.RUnlock()
				continue
			}
			// direct to client id
			if envelope.ToClient != "" {
				h.mu.RLock()
				for c := range h.clients {
					if c.ID == envelope.ToClient {
						select {
						case c.Send <- *envelope.Msg:
						default:
						}
					}
				}
				h.mu.RUnlock()
				continue
			}
			// broadcast to all
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.Send <- *envelope.Msg:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(room string, msg Message) {
	h.broadcast <- &MessageEnvelope{Msg: &msg, ToRoom: room}
}

// Convenience: publish a message to a client
func (h *Hub) SendToClient(clientID string, msg Message) {
	h.broadcast <- &MessageEnvelope{Msg: &msg, ToClient: clientID}
}

// handleIncoming handles messages that arrived from a client
func (h *Hub) handleIncoming(msg *Message, from *Client) {
	// Simple routing: if To is a room -> broadcast there, if To is user -> direct
	if msg.To != "" {
		// assume rooms prefixed by "room:" for clarity (not required)
		if _, ok := h.rooms[msg.To]; ok {
			h.BroadcastToRoom(msg.To, *msg)
			return
		}
		// else assume user id
		h.SendToClient(msg.To, *msg)
		return
	}
	// else broadcast to room of sender if set
	if from.Room != "" {
		h.BroadcastToRoom(from.Room, *msg)
		return
	}
	// fallback: broadcast to all
	h.broadcast <- &MessageEnvelope{Msg: msg}
}

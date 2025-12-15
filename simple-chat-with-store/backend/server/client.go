package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for sample — tighten in production
	},
}

// Client represents a single websocket connection
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan Message
	ID   string // user id
	Room string // current room
}

func NewClient(h *Hub, conn *websocket.Conn, id, room string) *Client {
	return &Client{Hub: h, Conn: conn, Send: make(chan Message, 256), ID: id, Room: room}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// attach sender if missing
		if msg.From == "" {
			msg.From = c.ID
		}
		c.Hub.handleIncoming(&msg, c)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS upgrades HTTP connection to a websocket and registers the client
func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	uid := q.Get("uid")
	room := q.Get("room")
	if uid == "" {
		http.Error(w, "uid required", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	client := NewClient(h, conn, uid, room)
	h.register <- client

	// Send join message into hub (optional)
	hubMsg := Message{Type: MsgJoin, From: uid, To: room, Timestamp: time.Now().Unix()}
	h.BroadcastToRoom(room, hubMsg)

	// Start pumps
	go client.WritePump()
	client.ReadPump()
}

package server

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub   *Hub
	Conn  *websocket.Conn
	Send  chan Message
	DocID string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			break
		}

		// Update document
		doc, ok := c.Hub.store.GetDocument(msg.DocID)
		if ok {
			doc.Blocks = append(doc.Blocks, msg.Block)
		}

		c.Hub.broadcast <- &msg
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		c.Conn.WriteJSON(msg)
	}
}

package server

type Hub struct {
	clients    map[*Client]bool
	docRooms   map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	store      *Store
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		docRooms:   make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
		store:      NewStore(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			if h.docRooms[c.DocID] == nil {
				h.docRooms[c.DocID] = make(map[*Client]bool)
			}
			h.docRooms[c.DocID][c] = true

		case c := <-h.unregister:
			delete(h.clients, c)
			delete(h.docRooms[c.DocID], c)
			close(c.Send)

		case msg := <-h.broadcast:
			for c := range h.docRooms[msg.DocID] {
				select {
				case c.Send <- *msg:
				default:
				}
			}
		}
	}
}

package server

type MessageType string

const (
	EditBlock          MessageType = "edit_block"
	AddBlock           MessageType = "add_block"
	CursorMove         MessageType = "cursor_move"
	Identity           MessageType = "identity"
	ClientDisconnected MessageType = "client_disconnected"
	SyncDocument       MessageType = "sync_document"
)

type Cursor struct {
	BlockID string `json:"blockId"`
	Offset  int    `json:"offset"`
}

type Message struct {
	Type     MessageType `json:"type"`
	DocID    string      `json:"docId"`
	ClientID string      `json:"clientId"`
	Block    Block       `json:"block,omitempty"`
	Blocks   []Block     `json:"blocks,omitempty"`
	Cursor   Cursor      `json:"cursor,omitempty"`
}

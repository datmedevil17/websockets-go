package server

type MessageType string

const (
	EditBlock MessageType = "edit_block"
	AddBlock  MessageType = "add_block"
)

type Message struct {
	Type  MessageType `json:"type"`
	DocID string      `json:"docId"`
	Block Block       `json:"block"`
}

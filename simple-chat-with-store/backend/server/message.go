package server

import "time"

type MessageType string

const (
	MsgChat   MessageType = "chat"
	MsgJoin   MessageType = "join"
	MsgLeave  MessageType = "leave"
	MsgSystem MessageType = "system"
)

type Message struct {
	Type      MessageType `json:"type"`
	Content   string      `json:"content,omitempty"`
	From      string      `json:"from,omitempty"`
	To        string      `json:"to,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

func NewChatMessage(from, content, to string) Message {
	return Message{
		Type:      MsgChat,
		From:      from,
		Content:   content,
		To:        to,
		Timestamp: time.Now().Unix(),
	}
}

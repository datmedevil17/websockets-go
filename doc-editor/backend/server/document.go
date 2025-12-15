package server

type BlockType string

const (
	TextBlock  BlockType = "text"
	ImageBlock BlockType = "image"
)

type Block struct {
	ID      string    `json:"id"`
	Type    BlockType `json:"type"`
	Content string    `json:"content"` // text OR base64 image
}

type Document struct {
	ID     string  `json:"id"`
	Blocks []Block `json:"blocks"`
}

package im

import "net"

type UserIdentify interface {
	String() string
}

type MessageType int

const (
	MessageTypePlainTxt = MessageType(iota)
	MessageTypeHTML
	MessageTypeJSON
	MessageTypeImage
	MessageTypeAudio
	MessageTypeVideo
	MessageTypeBlob
)

type Message struct {
	From    UserIdentify   `json:"from"`
	Unix    int64          `json:"unix"`
	Until   int64          `json:"until,omitempty"`
	Type    MessageType    `json:"type"`
	Content string         `json:"content,omitempty"`
	Ext     map[string]any `json:"ext,omitempty"`
}

type Channel struct {
	conns net.Conn
}

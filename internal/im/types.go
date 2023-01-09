package im

import "net"

type MessageType int

const (
	messageTypeMin = MessageType(iota)
	messageTypeUndefined
	MessageTypePlainTxt
	MessageTypeHTML
	MessageTypeJSON
	MessageTypeImage
	MessageTypeAudio
	MessageTypeVideo
	MessageTypeBlob
	messageTypeMax
)

type Message struct {
	From    uint64
	Unix    uint64
	Until   uint64
	Type    MessageType
	Content []byte
	Ext     map[string]string
}

type Channel struct {
	conns net.Conn
}

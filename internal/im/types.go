package im

import (
	"bytes"
	"compress/gzip"
	"github.com/zzztttkkk/ink/internal/utils"
	"io"
	"net"
)

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

// IsCompressed
// https://docs.fileformat.com/compression/gz/#gz-file-header
func (msg *Message) IsCompressed() bool {
	if len(msg.Content) < 10 {
		return false
	}
	return msg.Content[0] == 0x1f && msg.Content[1] == 0x8b
}

func (msg *Message) Compress() *bytes.Buffer {
	if len(msg.Content) < 512 {
		return nil
	}

	buf := utils.BytesBufferPool.Get()
	w := gzip.NewWriter(buf)
	_, e := w.Write(msg.Content)
	if e != nil {
		panic(e)
	}
	msg.Content = buf.Bytes()
	return buf
}

func (msg *Message) Uncompress() {
	raw := msg.Content
	r, e := gzip.NewReader(bytes.NewBuffer(raw))
	if e != nil {
		panic(e)
	}

	msg.Content = nil

	var buf [128]byte
	for {
		l, e := r.Read(buf[:])
		if e != nil {
			if e == io.EOF {
				break
			}
			panic(e)
		}
		if l < 1 {
			break
		}
		msg.Content = append(msg.Content, buf[:l]...)
	}
}

type Channel struct {
	conns net.Conn
}

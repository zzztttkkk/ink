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
	Ext     map[string]string
	Content []byte
}

// IsCompressed
// https://docs.fileformat.com/compression/gz/#gz-file-header
func (msg *Message) IsCompressed() bool {
	if len(msg.Content) < 10 {
		return false
	}
	return msg.Content[0] == 0x1f && msg.Content[1] == 0x8b
}

func (msg *Message) Compress(pptr **bytes.Buffer) {
	if len(msg.Content) < 512 {
		return
	}

	buf := utils.BytesBufferPool.Get()
	*pptr = buf

	w := gzip.NewWriter(buf)

	defer func() {
		if e := w.Close(); e != nil {
			panic(e)
		}
		msg.Content = buf.Bytes()
	}()

	_, _ = w.Write(msg.Content)
	return
}

func (msg *Message) Uncompress() {
	raw := msg.Content
	r, e := gzip.NewReader(bytes.NewBuffer(raw))
	if e != nil {
		panic(e)
	}
	defer r.Close()

	msg.Content = nil

	var buf [256]byte
	for {
		l, e := r.Read(buf[:])
		if e != nil && e != io.EOF {
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

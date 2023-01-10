package im

import (
	"bytes"
	"github.com/zzztttkkk/ink/internal/utils"
	"testing"
)

func TestMessage_Compress(t *testing.T) {
	msg := Message{
		From: 12,
	}

	for i := 0; i < 1024; i++ {
		msg.Content = append(msg.Content, "HX=="...)
	}

	rl := len(msg.Content)

	var buf *bytes.Buffer
	msg.Compress(&buf)
	defer utils.BytesBufferPool.Put(buf)

	msg.Uncompress()

	if len(msg.Content) != rl {
		t.Fail()
	}
}

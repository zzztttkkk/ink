package im

import (
	"fmt"
	"github.com/zzztttkkk/ink/internal/utils"
	"testing"
)

func TestA(t *testing.T) {
	msg := Message{
		From: 12,
	}

	for i := 0; i < 1024*1024; i++ {
		msg.Content = append(msg.Content, 'A')
	}

	buf := msg.Compress()
	defer utils.BytesBufferPool.Put(buf)

	msg.Uncompress()

	fmt.Println(1)
}

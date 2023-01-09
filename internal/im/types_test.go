package im

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
)

type UID int64

func (v UID) String() string {
	return strconv.FormatInt(int64(v), 16)
}

func TestA(t *testing.T) {
	msg := Message{
		From: UID(12),
		Unix: time.Now().Unix(),
		Text: "0.0",
		Ext:  nil,
	}
	v, _ := json.Marshal(msg)
	fmt.Println(string(v))
}

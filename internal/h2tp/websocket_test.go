package h2tp

import (
	"fmt"
	"github.com/zzztttkkk/ink/internal/utils"
	"github.com/zzztttkkk/ink/internal/vender/ws"
	"testing"
)

func TestWebsocket(t *testing.T) {
	router := NewRouter()

	router.RegisterWs("/ws", WsHandlerFunc(func(conn *WsConn) {
		for {
			data, opcode, e := conn.Read()
			if e != nil {
				break
			}

			if opcode != ws.OpText {
				continue
			}

			fmt.Println(utils.S(data))
			_ = conn.Write(opcode, data)
		}
	}))

	_ = Run("127.0.0.1:8000", map[string]*Router{"*": router})
}

package ink

import (
	"github.com/zzztttkkk/ink/internal/vender/ws"
	"github.com/zzztttkkk/ink/internal/vender/ws/wsutil"
	"net"
)

type WsServerSideConn struct {
	conn net.Conn
}

func (c *WsServerSideConn) Read() ([]byte, ws.OpCode, error) { return wsutil.ReadClientData(c.conn) }

func (c *WsServerSideConn) Write(op ws.OpCode, data []byte) error {
	return wsutil.WriteServerMessage(c.conn, op, data)
}

func Ws() Handler {
	return HandlerFunc(func(rctx *RequestCtx) {
		rctx.noTempResponse = true
		c, _, _, err := ws.UpgradeHTTP(rctx.Request, rctx.rw)
		if err != nil {
			return
		}
		defer c.Close()
	})
}

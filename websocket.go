package ink

import (
	"bufio"
	"github.com/zzztttkkk/ink/internal/vender/ws"
	"github.com/zzztttkkk/ink/internal/vender/ws/wsutil"
	"net"
)

type WsServerSideConn struct {
	Handshake ws.Handshake

	conn net.Conn
	rw   *bufio.ReadWriter
}

func (c *WsServerSideConn) Read() ([]byte, ws.OpCode, error) { return wsutil.ReadClientData(c.conn) }

func (c *WsServerSideConn) Write(op ws.OpCode, data []byte) error {
	return wsutil.WriteServerMessage(c.conn, op, data)
}

type WsHandler interface {
	Handle(conn *WsServerSideConn)
}

type WsHandlerFunc func(conn *WsServerSideConn)

func (fn WsHandlerFunc) Handle(conn *WsServerSideConn) { fn(conn) }

func makeWsHandler(handler WsHandler) Handler {
	return HandlerFunc(func(rctx *RequestCtx) {
		rctx.noTempResponse = true
		c, rw, hs, err := ws.UpgradeHTTP(rctx.Request, rctx.rw)
		if err != nil {
			return
		}
		defer c.Close()

		wsc := WsServerSideConn{conn: c, rw: rw, Handshake: hs}
		handler.Handle(&wsc)
	})
}

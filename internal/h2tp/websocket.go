package h2tp

import (
	"bufio"
	"github.com/zzztttkkk/ink/internal/vender/ws"
	"github.com/zzztttkkk/ink/internal/vender/ws/wsutil"
	"net"
)

type WsConn struct {
	Handshake ws.Handshake

	conn net.Conn
	rw   *bufio.ReadWriter
}

func (c *WsConn) Read() ([]byte, ws.OpCode, error) { return wsutil.ReadClientData(c.conn) }

func (c *WsConn) Write(op ws.OpCode, data []byte) error {
	return wsutil.WriteServerMessage(c.conn, op, data)
}

type WsHandler interface {
	Handle(conn *WsConn)
}

type WsHandlerFunc func(conn *WsConn)

func (fn WsHandlerFunc) Handle(conn *WsConn) { fn(conn) }

func makeWsHandler(handler WsHandler) Handler {
	return HandlerFunc(func(rctx *RequestCtx) {
		rctx.noTempResponse = true
		c, rw, hs, err := ws.UpgradeHTTP(rctx.Request, rctx.rw)
		if err != nil {
			return
		}
		defer c.Close()

		wsc := WsConn{conn: c, rw: rw, Handshake: hs}
		handler.Handle(&wsc)
	})
}

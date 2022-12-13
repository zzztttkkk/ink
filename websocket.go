package ink

import "nhooyr.io/websocket"

type WebSocketConn websocket.Conn
type WebSocketAcceptOptions websocket.AcceptOptions

type WebSocketHandler interface {
	Handle(c *WebSocketConn)
}

type WebSocketHandlerFunc func(*WebSocketConn)

func (w WebSocketHandlerFunc) Handle(c *WebSocketConn) { w(c) }

func Ws(handler WebSocketHandler, opts *WebSocketAcceptOptions) Handler {
	return HandlerFunc(func(rctx *RequestCtx) {
		rctx.noTempResponse = true
		c, err := websocket.Accept(rctx.rw, rctx.Request, (*websocket.AcceptOptions)(opts))
		if err != nil {
			return
		}
		defer c.Close(websocket.StatusNormalClosure, "")
		handler.Handle((*WebSocketConn)(c))
	})
}

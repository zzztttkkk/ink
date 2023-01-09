package h2tp

type Handler interface {
	Handle(rctx *RequestCtx)
}

type HandlerFunc func(rctx *RequestCtx)

func (f HandlerFunc) Handle(rctx *RequestCtx) { f(rctx) }

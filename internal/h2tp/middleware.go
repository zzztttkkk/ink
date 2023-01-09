package h2tp

type Middleware interface {
	Handle(rctx *RequestCtx, next func())
}

type MiddlewareFunc func(rctx *RequestCtx, next func())

func (f MiddlewareFunc) Handle(rctx *RequestCtx, next func()) { f(rctx, next) }

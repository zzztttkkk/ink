package ink

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/zzztttkkk/ink/internal/vld"
	"net/http"
	"reflect"
	"time"
)

type RequestCtx struct {
	Request    *http.Request
	PathParams httprouter.Params
	Response   Response

	noTempResponse bool
	rw             http.ResponseWriter
}

func (rctx *RequestCtx) Deadline() (deadline time.Time, ok bool) {
	return rctx.Request.Context().Deadline()
}

func (rctx *RequestCtx) Done() <-chan struct{} { return rctx.Request.Context().Done() }

func (rctx *RequestCtx) Err() error { return rctx.Request.Context().Err() }

func (rctx *RequestCtx) Value(key any) any { return rctx.Request.Context().Value(key) }

func (rctx *RequestCtx) NoTempResponse() bool { return rctx.noTempResponse }

var _ context.Context = (*RequestCtx)(nil)

func (rctx *RequestCtx) BindAndValidate(dist any) error {
	v := reflect.ValueOf(dist)
	return vld.GetRules(v.Type().Elem()).BindAndValidate(rctx.Request, v.Elem())
}

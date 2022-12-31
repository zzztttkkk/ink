package ink

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	router := NewRouter()

	router.Register(http.MethodGet, "/spk", HandlerFunc(func(rctx *RequestCtx) {
		rctx.Response.StatusCode = StatusOK
		_, _ = rctx.Response.WriteString("Hello World")
	}))

	router.Use(MiddlewareFunc(func(rctx *RequestCtx, next func()) {
		defer func() {
			fmt.Println("Middleware A After", time.Now().UnixNano())
		}()

		fmt.Println("Middleware A Before", time.Now().UnixNano(), fmt.Sprintf("%p", next))
		next()
	}))

	router.Use(MiddlewareFunc(func(rctx *RequestCtx, next func()) {
		defer func() {
			fmt.Println("Middleware B Defer After", time.Now().UnixNano())
		}()

		fmt.Println("Middleware B Before", time.Now().UnixNano(), fmt.Sprintf("%p", next))
		next()
		fmt.Println("Middleware B After", time.Now().UnixNano())
	}))

	group := NewGroup("/account")
	group.Use(MiddlewareFunc(func(rctx *RequestCtx, next func()) {
		fmt.Println("Middleware Group A")
		next()
	}))

	group.Register(http.MethodGet, "/register", HandlerFunc(func(rctx *RequestCtx) {
		_ = rctx.Response.WriteJSON(map[string]any{"xx": 34})
		panic(StatusBadRequest)
	}))

	router.AddGroup(group)

	if e := Run("127.0.0.1:8524", map[string]*Router{"*": router}); e != nil {
		panic(e)
	}
}

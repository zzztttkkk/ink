package ink

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/zzztttkkk/ink/internal/utils"
)

type RecoverFunc func(v any) Error

type Router struct {
	internal   *httprouter.Router
	middleware []Middleware

	Recover  RecoverFunc
	NotFound Handler
}

func NewRouter() *Router {
	obj := &Router{
		Recover: anyToError,
	}
	obj.internal = httprouter.New()
	return obj
}

func (r *Router) Use(middleware Middleware) {
	r.middleware = append(r.middleware, middleware)
}

func makeMiddlewareWrapper(middleware []Middleware, handler Handler) Handler {
	idx := -1

	return HandlerFunc(func(rctx *RequestCtx) {
		var next func()
		next = func() {
			idx++
			if idx < len(middleware) {
				middleware[idx].Handle(rctx, next)
			} else {
				handler.Handle(rctx)
			}
		}
		next()
	})
}

var AllMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

var (
	responseWBufPool = sync.Pool{New: func() any { return &bytes.Buffer{} }}
	MaxBufCap        = 10 << 20
)

func (r *Router) Register(methods string, pattern string, handler Handler) {
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	var temp []string
	if methods == "*" {
		temp = AllMethods
	} else {
		for _, part := range strings.Split(methods, ",") {
			part = strings.ToUpper(strings.TrimSpace(part))
			if utils.SliceFind(AllMethods, part) < 0 {
				panic(fmt.Errorf("unknown method, %s", part))
			}
			temp = append(temp, part)
		}
	}

	for _, method := range temp {
		r.internal.Handle(method, pattern, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			rctx := RequestCtx{
				Request:    request,
				PathParams: params,
				rw:         writer,
			}

			defer func() {
				defer func() {
					buf := rctx.Response.buf
					if buf == nil {
						return
					}
					if buf.Cap() > MaxBufCap {
						return
					}
					buf.Reset()
					responseWBufPool.Put(buf)
				}()

				ev := recover()
				if ev != nil && (rctx.Response.StatusCode == 0 || rctx.Response.StatusCode == StatusOK) {
					rctx.noTempResponse = false

					rctx.Response.ResetAll()
					err := r.Recover(ev)
					rctx.Response.StatusCode = err.StatusCode()
					rctx.Response.header = err.Header()
					_, _ = rctx.Response.Write(err.Body())
				}

				if !rctx.noTempResponse {
					for k, vs := range rctx.Response.header {
						fmt.Println(k, vs)
						writer.Header()[k] = vs
					}
					if rctx.Response.StatusCode == 0 {
						rctx.Response.StatusCode = StatusOK
					}

					writer.WriteHeader(rctx.Response.StatusCode)
					if rctx.Response.buf != nil {
						_, _ = writer.Write(rctx.Response.buf.Bytes())
					}
				}
			}()

			handler = makeMiddlewareWrapper(r.middleware, handler)
			handler.Handle(&rctx)
		})
	}
}

func (r *Router) RegisterFs(methods string, prefix string, filesystem any, opts *FsOpts) {
	if !strings.HasSuffix(prefix, "/") {
		panic(fmt.Sprintf("bad prefix, `%s`", prefix))
	}
	prefix += fmt.Sprintf("*%s", filepathArgName)
	var optsVal FsOpts
	if opts != nil {
		optsVal = *opts
	}
	r.Register(methods, prefix, makeFsHandler(filesystem, optsVal))
}

func (r *Router) AddGroup(group *Group) {
	for _, info := range group.infos {
		methods, pattern, handler := info.methods, info.pattern, info.handler
		handler = makeMiddlewareWrapper(group.middleware, handler)
		r.Register(methods, pattern, handler)
	}
}

func Run(addr string, routers map[string]*Router) error {
	if len(routers) < 1 {
		routers = map[string]*Router{}
		router := NewRouter()
		routers["*"] = router
	}

	defaultRouter := routers["*"]
	delete(routers, "*")

	var peekRouter func(r *http.Request) *Router
	if len(routers) == 0 {
		peekRouter = func(_ *http.Request) *Router {
			return defaultRouter
		}
	} else {
		if len(routers) == 1 {
			host := utils.MapKeys(routers)[0]
			router := utils.MapValues(routers)[0]
			peekRouter = func(req *http.Request) *Router {
				if req.Host != host {
					return nil
				}
				return router
			}
		} else {
			if len(routers) < 6 {
				_hosts := utils.MapKeys(routers)
				_routers := utils.MapValues(routers)
				peekRouter = func(r *http.Request) *Router {
					for i, host := range _hosts {
						if r.Host == host {
							return _routers[i]
						}
					}
					return nil
				}
			} else {
				peekRouter = func(r *http.Request) *Router {
					return routers[r.Host]
				}
			}
		}
	}

	return http.ListenAndServe(addr, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		router := peekRouter(request)
		if router == nil {
			router = defaultRouter
		}

		if router == nil {
			return
		}

		fn, params, _ := router.internal.Lookup(request.Method, request.RequestURI)
		if fn == nil {
			if router.NotFound == nil {
				writer.WriteHeader(StatusNotFound)
				return
			}
			router.NotFound.Handle(&RequestCtx{Request: request, PathParams: params, rw: writer})
			return
		}

		fn(writer, request, params)
	}))
}

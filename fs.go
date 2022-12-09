package h2tp

import (
	"fmt"
	fslib "io/fs"
	"net/http"
	"reflect"
	"strings"
)

type OsDirOpts struct {
	URLPrefix    string
	DisableIndex bool
}

var (
	ioFsType   = reflect.TypeOf((*fslib.FS)(nil)).Elem()
	httpFsType = reflect.TypeOf((*http.FileSystem)(nil)).Elem()
)

func Fs(fs any, opt OsDirOpts) Handler {
	raw := reflect.ValueOf(fs).Interface()

	fsv := reflect.ValueOf(fs)
	if fsv.Kind() == reflect.String {
		fsv = reflect.ValueOf(http.Dir(raw.(string)))
	}

	if fsv.Type().Implements(ioFsType) {
		fsv = reflect.ValueOf(http.FS(fsv.Interface().(fslib.FS)))
	}

	if !fsv.Type().Implements(httpFsType) {
		panic(fmt.Errorf("%v can not cast to http.FileSystem", raw))
	}

	httpfs := http.FileServer(fsv.Interface().(http.FileSystem))
	return HandlerFunc(func(rctx *RequestCtx) {
		rctx.noTempResponse = true
		if opt.DisableIndex && strings.HasSuffix(rctx.Request.URL.Path, "/") {
			rctx.rw.WriteHeader(StatusMovedPermanently)
			return
		}
		if len(opt.URLPrefix) > 0 && strings.HasPrefix(rctx.Request.URL.Path, opt.URLPrefix) {
			rctx.Request.URL.Path = rctx.Request.URL.Path[len(opt.URLPrefix):]
		}
		httpfs.ServeHTTP(rctx.rw, rctx.Request)
	})
}

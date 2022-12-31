package ink

import (
	"fmt"
	"io/fs"
	"net/http"
	"reflect"
	"strings"
)

type FsOpts struct {
	DisableIndex  bool
	IndexRenderer func(rctx *RequestCtx, infos []fs.FileInfo)
}

var (
	ioFsType        = reflect.TypeOf((*fs.FS)(nil)).Elem()
	httpFsType      = reflect.TypeOf((*http.FileSystem)(nil)).Elem()
	filepathArgName = "filepath"
)

func makeFsHandler(anyFs any, opt FsOpts) Handler {
	raw := reflect.ValueOf(anyFs).Interface()

	fsv := reflect.ValueOf(anyFs)
	if fsv.Kind() == reflect.String {
		fsv = reflect.ValueOf(http.Dir(raw.(string)))
	}

	if fsv.Type().Implements(ioFsType) {
		fsv = reflect.ValueOf(http.FS(fsv.Interface().(fs.FS)))
	}

	if !fsv.Type().Implements(httpFsType) {
		panic(fmt.Errorf("%v can not cast to http.FileSystem", raw))
	}

	httpfs := fsv.Interface().(http.FileSystem)
	server := http.FileServer(httpfs)
	return HandlerFunc(func(rctx *RequestCtx) {
		rctx.Request.URL.Path = rctx.PathParams.ByName(filepathArgName)
		maybeDir := strings.HasSuffix(rctx.Request.URL.Path, "/")
		if opt.DisableIndex && maybeDir {
			rctx.rw.WriteHeader(StatusMovedPermanently)
			return
		}
		if maybeDir && opt.IndexRenderer != nil {
			f, e := httpfs.Open(rctx.Request.URL.Path)
			if e != nil {
				panic(e)
			}
			s, e := f.Stat()
			if e != nil {
				panic(e)
			}
			if s.IsDir() {
				infos, err := f.Readdir(-1)
				if err != nil {
					panic(err)
				}
				opt.IndexRenderer(rctx, infos)
				return
			}
		}
		rctx.noTempResponse = true
		server.ServeHTTP(rctx.rw, rctx.Request)
	})
}

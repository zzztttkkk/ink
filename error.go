package ink

import (
	"github.com/zzztttkkk/ink/internal/utils"
	"net/http"
	"reflect"
)

type Error interface {
	StatusCode() int
	Header() http.Header
	Body() []byte
}

var (
	errInterfaceType = reflect.TypeOf((*Error)(nil)).Elem()
)

type StatusError int

func (status StatusError) StatusCode() int { return int(status) }

func (status StatusError) Header() http.Header { return nil }

func (status StatusError) Body() []byte { return utils.B(http.StatusText(int(status))) }

var _ Error = (StatusError)(0)

func anyToError(ev any) Error {
	vv := reflect.ValueOf(ev)
	vt := vv.Type()

	if vt.Kind() == reflect.Int {
		return StatusError(ev.(int))
	}

	if vt.Implements(errInterfaceType) {
		return ev.(Error)
	}
	return StatusError(StatusInternalServerError)
}

package vld

import "net/http"

type Binder interface {
	FromRequest(*http.Request) *Error
}

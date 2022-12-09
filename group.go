package ink

import "strings"

type _LazyRegisterInfo struct {
	methods string
	pattern string
	handler Handler
}

type Group struct {
	middleware []Middleware
	prefix     string
	infos      []_LazyRegisterInfo
}

func NewGroup(prefix string) *Group {
	return &Group{prefix: prefix}
}

func (group *Group) Use(middleware Middleware) {
	group.middleware = append(group.middleware, middleware)
}

func (group *Group) Register(methods string, pattern string, handler Handler) {
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	group.infos = append(group.infos, _LazyRegisterInfo{methods, group.prefix + pattern, handler})
}

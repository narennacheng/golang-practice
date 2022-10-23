package geek

import (
	"fmt"
)

type Routable interface {
	Route(method, pattern string, handlerFunc func(ctx *Context)) error
}

// Handler 基于原生http.Handler 扩展自己功能
type Handler interface {
	ServeHTTP(ctx *Context)
	Routable
}

// HandlerBasedOnMap 基于 map 的路由
type HandlerBasedOnMap struct {
	// key: method + urlPath
	handlers map[string]func(ctx *Context)
}

// Route 设定一个路由，命中该路由的会执行handlerFunc
func (h *HandlerBasedOnMap) Route(method, pattern string, handlerFunc func(ctx *Context)) error {
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
	return nil
}

func (h *HandlerBasedOnMap) ServeHTTP(ctx *Context) {
	key := h.key(ctx.R.Method, ctx.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(ctx)
	} else {
		ctx.NotFound404()
	}
}

func (h *HandlerBasedOnMap) key(method, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

// 一种常用的Go设计模式，静态检查 确保HandlerBasedOnMap实现了Handler
var _ Handler = &HandlerBasedOnMap{}

func NewHandlerBaseOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(ctx *Context), 10),
	}
}

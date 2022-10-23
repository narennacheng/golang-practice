package geek

import (
	"net/http"
	"testing"
)

/*
golang 单元测试规范：
1.文件用xxx_test.go结尾
2.方法形式 TestXXXX(t *testing.T)

一般来说，单元测试都要：
1. 正常执行
2. 异常例子
3. 正常与异常的边界例子
4. 追求分支覆盖(重点)
5. 追求代码覆盖
*/

func TestHandlerBasedOnTree_Route(t *testing.T) {
	handler := NewHandlerBasedOnTree().(*HandlerBasedOnTree)
	handler.Route(http.MethodPost, "/user/info", func(ctx *Context) {})
	PrintChild(handler.root)

	handler.Route(http.MethodGet, "/user", func(ctx *Context) {})

	PrintChild(handler.root)

	handler.Route(http.MethodGet, "/user/info/name", func(ctx *Context) {})
	handler.Route(http.MethodGet, "/user/login", func(ctx *Context) {})
	PrintChild(handler.root)

	handler.Route(http.MethodGet, "/order", func(ctx *Context) {})
	handler.Route(http.MethodGet, "/order/user", func(ctx *Context) {})
	PrintChild(handler.root)

}

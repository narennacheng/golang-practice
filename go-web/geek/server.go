package geek

import (
	ccontext "context"
	"fmt"
	"net/http"
)

type Server interface {
	// Route 设定一个路由，命中该路由的会执行handlerFunc
	Routable
	// Start 启动server
	Start(address string) error
	// Shutdown 关闭server
	Shutdown(ctx ccontext.Context) error
}

// sdkHttpServer 基于 http 实现
type sdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

// Route 设定一个路由，命中该路由的会执行handlerFunc
func (s *sdkHttpServer) Route(method, pattern string, handlerFunc func(ctx *Context)) error {
	return s.handler.Route(method, pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		s.root(c)
	})
	fmt.Printf("server[%v] listening %v...\n", s.Name, address)
	return http.ListenAndServe(address, nil)
}

func (s *sdkHttpServer) Shutdown(ctx ccontext.Context) error {
	return nil
}

// NewSdkHttpServer 定义完接口后，通常暴露一个创建方法
func NewSdkHttpServer(name string, builders ...FilterBuiler) Server {
	// handler := NewHandlerBaseOnMap()
	handler := NewHandlerBasedOnTree()
	// 因为是一个链，所以要把最后的业务逻辑处理
	var root Filter = handler.ServeHTTP
	// 从后往前把filter串起来
	for i := len(builders) - 1; i >= 0; i-- {
		root = builders[i](root)
	}
	return &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}

type signUpReq struct {
	// Tag
	Email    string `json:"email"`
	Password string `json:"password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Signup(ctx *Context) {
	req := &signUpReq{}
	err := ctx.ReadJson(req)
	if err != nil {
		_ = ctx.BadRequestJson(&commonResponse{
			BizCode: 4,
			// 实际中，应避免暴露 error
			Msg: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}
	_ = ctx.OkJson(
		// 登录成功，返回用户id。
		&commonResponse{
			Data: 1,
		},
	)

}

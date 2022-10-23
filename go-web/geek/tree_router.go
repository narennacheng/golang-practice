package geek

import (
	"errors"
	"fmt"
	"strings"
)

/*
路由树 定义
*/

type node struct {
	path     string
	children []*node

	// 如果这个叶子节点，那么匹配上后就调用该方法
	handler handlerFunc
}

func PrintChild(n *node) {
	if len(n.children) == 0 {
		fmt.Printf("path:/%v, no have child\n", n.path)
		return
	}
	for i, nc := range n.children {
		fmt.Printf("path:/%v, child[%d]: %+v\n", n.path, i, nc)
		PrintChild(nc)
	}
}

func (n *node) findMatchChild(path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range n.children {
		if child.path == path && child.path != "*" {
			return child, true
		}
		// 命中通配符
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
}

func (n *node) createSubTree(paths []string, handler handlerFunc) {
	for _, path := range paths {
		nn := newNode(path)
		n.children = append(n.children, nn)
		n = nn
	}
	n.handler = handler
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 2),
	}
}

// HandlerBasedOnTree 基于 tree 的路由
type HandlerBasedOnTree struct {
	root *node
}

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")

// validatePattern 校验路由书写规则
func (h *HandlerBasedOnTree) validatePattern(pattern string) error {
	// 校验 *，如果存在，必须在最后一个，并且它前面必须是/
	// 即我们只接受 /* 的存在
	pos := strings.Index(pattern, "*")
	if pos > 0 {
		if pos != len(pattern)-1 {
			return ErrorInvalidRouterPattern
		}
		if pattern[pos-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return nil
}

func (h *HandlerBasedOnTree) Route(method, pattern string, handlerFunc func(ctx *Context)) error {
	err := h.validatePattern(pattern)
	if err != nil {
		return err
	}

	// 去除前后'/'，再分割；/user/info/ => user/info => [user, info]
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")
	// 指向当前节点
	cur := h.root
	for index, path := range paths {
		matchChild, found := cur.findMatchChild(path)
		if found {
			cur = matchChild
		} else {
			cur.createSubTree(paths[index:], handlerFunc)
			return nil
		}
	}
	// 离开了循环说明加入的是短路径
	// 比如先加入 /order/detail 再加入/order
	cur.handler = handlerFunc
	return nil
}

func (h *HandlerBasedOnTree) ServeHTTP(ctx *Context) {
	url := strings.Trim(ctx.R.URL.Path, "/")
	paths := strings.Split(url, "/")
	cur := h.root
	for _, path := range paths {
		// 从子节点里找一个匹配到了当前 path 的节点
		matchChild, found := cur.findMatchChild(path)
		if !found {
			// 如果找不到，直接返回
			ctx.NotFound404()
			return
		}
		cur = matchChild
	}
	// 到这里应该是遍历找完了
	if cur.handler == nil {
		// 到达这里是因为 比如注册了/user/info 然后访问/user
		ctx.NotFound404()
		return
	}
	cur.handler(ctx)
}

// 一种常用的Go设计模式，静态检查 确保HandlerBasedOnTree实现了Handler
var _ Handler = &HandlerBasedOnTree{}

func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: &node{},
	}
}

package main

import (
	"context"
	"fmt"
	"time"
)

/*
context.Context 是 Go 提供的线程安全工具，称为上下文。
方法：
	• WithTimeout：一般用户控制超时
	• WithCancel：用于取消整条链上的任务
	• WithDeadline：控制时间
	• WithValue：往里面塞入 key-value
	• Backgroud：返回一个空的 context.Context
	• ToDo：返回一个空的 context.Context，但是这个标记着你也不知道传什么
*/
func main() {
	// withTimeout()
	// withDeadline()
	withValue()
}

// WithTimeout：一般用户控制超时
func withTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	start := time.Now().Unix()
	<-ctx.Done()
	end := time.Now().Unix()
	fmt.Println("withTimeout: ", end-start) // 输出2。说明在ctx.Done等待了2s
}

// WithDeadline：控制时间
func withDeadline() {
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cf()

	start := time.Now().Unix()
	<-ctx.Done()
	end := time.Now().Unix()
	fmt.Println("withDeadline: ", end-start) // 输出2。说明在ctx.Done等待了2s

}

// WithValue：往里面塞入 key-value
func withValue() {
	parentKey := "parent"
	ctx := context.WithValue(context.Background(), parentKey, "this is parent")

	sonKey := "son"
	ctx2 := context.WithValue(ctx, sonKey, "this is son")

	// 尝试从 parent 里面拿出来 sonKey，会拿不到，只能取出 parentKey
	if ctx.Value(sonKey) == nil {
		fmt.Println("parent can not get son's key-value pair")
	}
	if val := ctx2.Value(parentKey); val != nil {
		fmt.Printf("son can get parent's key-value: %v", val)
	}
}

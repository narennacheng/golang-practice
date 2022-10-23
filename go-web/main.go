package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"geek"
)

func main() {
	shutdown := geek.NewGracefulShutdown()
	server := geek.NewSdkHttpServer("test-server",
		geek.MetricFilterBuiler, shutdown.ShutdownFilterBuilder)
	// if err := server.Start(":8080"); err != nil {
	// 	panic(err)
	// }

	adminServer := geek.NewSdkHttpServer("admin-server",
		// 注意：实际生产中，使用多个 server 监听不同端口，
		// 那么shutdown 最好也是多个，互相之间就不有竞争
		geek.MetricFilterBuiler, shutdown.ShutdownFilterBuilder)

	// 注册路由
	_ = server.Route(http.MethodPost, "/user/login", geek.Signup)

	go func() {
		if err := adminServer.Start(":8081"); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := server.Start(":8080"); err != nil {
			// 假设还有其他操作
			panic(err)
		}
	}()

	// WaitForShutdown监听等待信号，执行完后才退出main函数
	// 先执行RejectNewRequestAndWaiting，等待所有请求
	// 然后关闭 server，如果是多个server，可以多个goroutine一起关闭
	geek.WaitForShutdown(
		// 假设这里有个 hook 通知网关我们下线了
		func(ctx context.Context) error {
			fmt.Println("mock notify gateway")
			time.Sleep(time.Second * 2)
			return nil
		},
		// RejectNewRequestAndWaiting hook
		shutdown.RejectNewRequestAndWaiting,
		// 全部请求处理完就可以关闭 server hook
		geek.BuildCloseServerHook(server, adminServer),
		// 假设这里有个 hook 用来关闭临时资源
		func(ctx context.Context) error {
			fmt.Println("mock release resources")
			time.Sleep(time.Second * 2)
			return nil
		},
	)
}

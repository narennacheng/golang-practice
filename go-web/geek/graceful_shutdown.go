package geek

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var ErrorHookTimeout = errors.New("hook timeout")

// 等待当前的所有请求处理完毕：需要维持请求计数
type GracefulShutdown struct {
	// 还在处理中的请求数
	reqCnt int64
	// 大于 1 说明要关闭
	closing int32

	// 用 channel 通知已经处理完所有请求
	zeroReqCnt chan struct{}
}

// RejectNewRequestAndWaiting 拒绝新请求，并等待正在处理中的请求
func (g *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {
	atomic.AddInt32(&g.closing, 1)

	// 如果关闭前就已经处理完请求，直接返回
	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}
	done := ctx.Done()
	// 因为是单向的，所以这里不用for
	// 单向：一触发就回不到原理正常处理请求的状态了
	// 这个 select 可以理解为 要么超时，要么这里所有请求都执行完了
	select {
	case <-done:
		fmt.Println("超时了，还没等到所有请求执行完毕...")
		return ErrorHookTimeout
	case <-g.zeroReqCnt:
		fmt.Println("全部请求处理完毕...")
	}
	return nil
}

func (g *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		cl := atomic.LoadInt32(&g.closing)
		// 已经接收到关闭信号, 不接受新的请求
		if cl > 0 {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// 维持请求计数
		atomic.AddInt64(&g.reqCnt, 1)
		next(c)
		n := atomic.AddInt64(&g.reqCnt, -1)

		// 已经开始关闭，而且请求数为0
		if cl > 0 && 0 == n {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

func WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdowmSignalsWin...)
	select {
	case sig := <-signals:
		fmt.Printf("get signal:%s, application will shutdown \n", sig)
		// 十分钟后退出
		time.AfterFunc(time.Minute*10, func() {
			fmt.Printf("shutdown gracefully timeout, application")
			os.Exit(1)
		})
		for _, h := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			err := h(ctx)
			if err != nil {
				fmt.Printf("failed to return hook, err:%v\n", err)
			}
			cancel()
		}
		os.Exit(0)
	}
}

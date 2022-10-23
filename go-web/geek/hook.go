package geek

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 处理超时的问题的钩子
type Hook func(ctx context.Context) error

func BuildCloseServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{} // goroutine 之间的同步
		doneCh := make(chan struct{})
		// 可能传了多个server
		wg.Add(len(servers))

		for _, s := range servers {
			go func(svr Server) {
				err := svr.Shutdown(ctx)
				if err != nil {
					fmt.Printf("server shutdown error: %v \n", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(s)
		}

		go func() {
			// wait会阻塞，直到Done()调用使得计数归为0
			wg.Wait()
			doneCh <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("closing servers timeout \n")
			return ErrorHookTimeout
		case <-doneCh:
			fmt.Println("close all servers \n")
			return nil
		}

	}
}

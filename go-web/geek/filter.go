package geek

import (
	"fmt"
	"time"
)

type handlerFunc func(c *Context)

type FilterBuiler func(next Filter) Filter

type Filter func(c *Context)

// MetricFilterBuiler 定义一个打印执行时间的filter
func MetricFilterBuiler(next Filter) Filter {
	return func(c *Context) {
		startTime := time.Now().UnixNano()
		next(c)
		// time.Sleep(time.Second)
		endTime := time.Now().UnixNano()
		fmt.Printf("used time: %dns \n", endTime-startTime)
	}
}

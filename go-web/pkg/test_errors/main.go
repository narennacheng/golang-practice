package main

import (
	"errors"
	"fmt"
)

/*
errors 包
• New 创建一个新的 error
• Is 判断是不是特定的某个error
• As 类型转换为特定的error
• Unwrap 解除包装，返回被包装的 error
*/

type MyError struct {
}

func (m *MyError) Error() string {
	return "its my error"
}

func ErrorsPkg() {
	err := &MyError{}
	// 使用 %w 占位符，返回的是一个新错误
	// wrappedErr 是一个新类型，
	wrappedErr := fmt.Errorf("this is an wrapped error %w\n", err)

	// 再解出来
	if err == errors.Unwrap(wrappedErr) {
		println("unwrapped")
	}

	if errors.Is(wrappedErr, err) {
		// 虽然包了一下，但是 Is 会逐层解除包装，判断是不是该错误
		println("wrapped is err")
	}

	err2 := &MyError{}
	// 尝试将 wrappedErr 转换为 MyError
	if errors.As(wrappedErr, &err2) {
		println("convert error")
	}
}

func main() {
	ErrorsPkg()

	println("---测试 panic 后 recover --")

	defer func() {
		if data := recover(); data != nil {
			fmt.Printf("hello, panic:%v \n", data)
		}
		println("恢复后从这里继续执行")
	}()
	panic("boom!!")
	println("这里不会执行")
}

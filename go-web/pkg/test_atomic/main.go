package main

import "sync/atomic"

/*
atomic 包
方法分成这几类：
	• AddXXX：操作一个数字类型，加上一个数字
	• LoadXXX：读取一个值
	• CompareAndSwapXXX：大名鼎鼎的 CAS 操作
	• StoreXXX：写入一个值
	• SwapXXX：写入一个值，并且返回旧的值。它和 CompareAndSwap 的区别在于它不关心旧的值是什么
	• unsafepointer 相关方法，不建议使用。
		难写也难读，不到逼不得已不要去用。尤其是不要为了优化而故意用 unsafepoint
*/

var value int32 = 0

func main() {
	// 要传入 value 的指针
	// value + 10
	atomic.AddInt32(&value, 10)
	nv := atomic.LoadInt32(&value)
	println(nv) // 10

	// 如果之前的值是10，那么设置为新的值20
	swapped := atomic.CompareAndSwapInt32(&value, 10, 20)
	println(swapped) // true

	old := atomic.SwapInt32(&value, 40)
	println(old)   // 20
	println(value) // 40
}

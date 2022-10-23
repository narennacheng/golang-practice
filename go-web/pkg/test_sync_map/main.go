package main

import (
	"fmt"
	"sync"
)

/*
类型断言: t, ok := x.(T) 或者 t :=x.(T)   // T 可以是结构体或者指针
类型转换: y := T(x)  // 数字类型转换，string 和 []byte 互相转
*/
func testType(val interface{}) {
	fmt.Println(len(val.(string)))
	a := 12.0
	b := int(a)
	println(b)

	str := "hello"
	bytes := ([]byte)(str)
	println(bytes)
}

/*
sync 包提供了基本的并发工具
• sync.Map：并发安全 map
• sync.Mutex：锁
• sync.RWMutex：读写锁
• sync.Once：只执行一次
• sync.WaitGroup: goroutine 之间同步

尽量用 RWMutext
• 尽量用 defer 来释放锁，防止panic没有释放锁
• 不可重入：lock 之后，即便是同一个线程(goroutine)，也无法再次加锁（写递归函数要小心）
• 不可升级：加了读锁之后，如果试图加写锁，锁不升级
*/
func testSyncMap() {
	// 线程安全的 Map
	m := sync.Map{}
	m.Store("key", "value")
	m.Store("key1", "value1")

	val, ok := m.Load("key1")
	if ok {
		testType(val)
	}
}

var mutex sync.Mutex
var rwMutex sync.RWMutex

func testMutex() {
	mutex.Lock()
	defer mutex.Unlock()
	// your code
}

func testRwMutex() {
	// 加读锁
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	// 加写锁
	rwMutex.Lock()
	defer rwMutex.Unlock()
}

// 不可重入例子
func failed1() {
	mutex.Lock()
	defer mutex.Unlock()

	// 这一句会死锁
	// 但是如果只有一个goroutine，那么这个会导致程序崩溃
	mutex.Lock()
	defer mutex.Unlock()
}

// 不可升级例子
func failed2() {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	// 这一句会死锁
	// 但是如果只有一个goroutine，那么这个会导致程序崩溃
	mutex.Lock()
	defer mutex.Unlock()
}

var once sync.Once

func testOnce() {
	once.Do(func() {
		fmt.Println("只执行一次")
	})
}

func testWaitGroup() {
	res := 0
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(val int) {
			res += val
			wg.Done()
		}(i)
	}
	wg.Wait() // 把这个注释 res可能就不是是45
	fmt.Println(res)
}

/*
1. 尽量用 sync.RWMutex
2. sync.Once 可以保证代码只会执行一次，一般用来解决一些初始化的需求
3. sync.WaitGroup 能用来在多个 goroutine 之间进行同步
*/
func main() {
	testSyncMap()
	testOnce()
	testWaitGroup()

}

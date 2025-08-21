package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var counter int
var mutex sync.Mutex

// 1增加计数器的值
func increment() {
	for i := 0; i < 1000; i++ {
		//加锁
		mutex.Lock()
		counter++
		//解锁
		mutex.Unlock()
	}
}

// 使用原子操作的计数器
var counter2 int32

// 2增加计数器的值
func increment2() {
	for i := 0; i < 1000; i++ {
		//原子递增计数器
		atomic.AddInt32(&counter2, 1)
	}
}

func main() {
	fmt.Println("------------Task1 Start------------")
	var wg sync.WaitGroup
	//启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		fmt.Println("协程启动:", i)
		go func() {
			defer wg.Done()
			increment()
		}()
	}
	//End
	wg.Wait()
	fmt.Println("all counter price:", counter)
	fmt.Println("------------Task1 End------------")

	fmt.Println("------------Task2 Start------------")
	var wg2 sync.WaitGroup
	//启动10个协程
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		fmt.Println("协程启动:", i)
		go func() {
			defer wg2.Done()
			increment2()
		}()
	}
	//End
	wg2.Wait()
	fmt.Println("all counter2 price:", counter2)
	fmt.Println("------------Task2 End------------")
}

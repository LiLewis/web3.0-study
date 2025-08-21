package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("------------Task1 Start------------")
	// 无缓冲通道
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)

	//生产者协程（1~10）
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
		//循环完毕关闭通道
		close(ch)
	}()

	//消费者协程
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("消费者接收: %d\n ", num)
		}
	}()

	//End
	wg.Wait()
	fmt.Println("------------Task1 End------------")

	fmt.Println("------------Task2 Start------------")
	// 带缓冲通道，容量:10
	ch2 := make(chan int, 10)
	var wg2 sync.WaitGroup

	wg2.Add(2)

	//生产者协程（1~100）
	go func() {
		defer wg2.Done()
		for i := 0; i < 100; i++ {
			ch2 <- i
		}
		//循环完毕关闭通道
		close(ch2)
	}()

	//消费者协程
	go func() {
		defer wg2.Done()
		for num := range ch2 {
			fmt.Printf("消费者接收: %d\n ", num)
		}
	}()

	//End
	wg2.Wait()
	fmt.Println("------------Task2 End------------")
}

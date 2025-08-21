package main

import (
	"fmt"
	"sync"
	"time"
)

// 奇数函数
func printOdd() {
	for i := 1; i <= 10; i += 2 {
		fmt.Println("打印奇数值：", i)
		//每次循环一次增加100毫秒
		time.Sleep(100 * time.Millisecond)
	}
}

// 偶数函数
func printEven() {
	for i := 2; i <= 10; i += 2 {
		fmt.Println("打印偶数值：", i)
		//每次循环一次增加100毫秒
		time.Sleep(100 * time.Millisecond)
	}
}

// Task 定义任务器
type Task func()

// 执行任务并统计耗时
func runTimesTask(id int, t Task, wg *sync.WaitGroup) {
	// goroutine 结束时计数
	defer wg.Done()
	start := time.Now()

	//执行
	t()

	elapsed := time.Since(start)
	fmt.Printf("任务 %d 执行完成，耗时: %v\n", id, elapsed)

}

func main() {
	fmt.Println("------------Task1 Start------------")
	//根据函数位置先后执行顺序进行执行
	printOdd()
	printEven()

	time.Sleep(2 * time.Second)
	fmt.Println("------------Task1 End------------")

	fmt.Println("------------Task2 Start------------")
	// Task 定义一组任务
	tasks := []Task{
		func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Task1: API接口请求")
		},
		func() {
			time.Sleep(800 * time.Millisecond)
			fmt.Println("Task2: Image加载......")
		},
		func() {
			time.Sleep(300 * time.Millisecond)
			fmt.Println("Task3: 数据传输")
		},
	}

	var wg sync.WaitGroup
	//并发处理
	for i, task := range tasks {
		wg.Add(1)
		go runTimesTask(i+1, task, &wg)
	}

	//End
	wg.Wait()
	fmt.Println("all done!!!!")
	fmt.Println("------------Task2 End------------")

}

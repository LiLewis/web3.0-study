package main

import "fmt"

func addPointerNums(num *int) {
	*num += 10
}

func sliceByDouble(arr *[]int) {
	if arr == nil {
		return
	}
	for i := range *arr {
		(*arr)[i] = (*arr)[i] * 2
	}
}

func main() {
	//执行addPointerNums函数
	oldNums := 5
	fmt.Println("入参：", oldNums)

	addPointerNums(&oldNums)

	newNums := oldNums
	fmt.Println("函数执行后的出参：", newNums)

	//执行sliceByDouble含糊
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println("入参：", arr)

	sliceByDouble(&arr)

	doubleArr := arr
	fmt.Println(doubleArr)

}

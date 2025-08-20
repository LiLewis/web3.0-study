package main

import "fmt"

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for k, num := range nums {
		complement := target - num
		if idx, ok := m[complement]; ok {
			return []int{idx, k}
		}
		m[num] = k
	}
	return nil
}

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
	fmt.Println(twoSum([]int{3, 2, 4}, 6))
	fmt.Println(twoSum([]int{3, 3}, 6))
}

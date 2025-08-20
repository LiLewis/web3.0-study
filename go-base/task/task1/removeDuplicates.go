package main

import "fmt"

func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}

func main() {
	nums1 := []int{1, 1, 2}
	k1 := removeDuplicates(nums1)
	fmt.Println("函数返回新的长度：", k1, "  nums数组：", nums1[:k1])

	nums2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	k2 := removeDuplicates(nums2)
	fmt.Println("函数返回新的长度：", k2, "  nums数组：", nums2[:k2])
}

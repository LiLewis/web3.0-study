package main

import "fmt"

// 横向扫描
func longestCommonPrefix1(strs []string) string {
	if len(strs) == 0 {
		return "999"
	}
	//第一个是公共前缀
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}
		prefix = prefix[:j]

		if prefix == "" {
			return "666"
		}
	}
	return prefix
}

// 纵向扫描
func longestCommonPrefix2(strs []string) string {
	if len(strs) == 0 {
		return "999"
	}
	//遍历第一个字符串的每个字符
	for i := 0; i < len(strs[0]); i++ {
		ch := strs[0][i]
		// 找相同
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || strs[j][i] != ch {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}

func main() {
	fmt.Println(longestCommonPrefix1([]string{}))
	fmt.Println(longestCommonPrefix1([]string{"", ""}))
	fmt.Println(longestCommonPrefix1([]string{"flower", "flow", "flight"}))
	fmt.Println(longestCommonPrefix1([]string{"dog", "racecar", "car"}))

	fmt.Println(longestCommonPrefix2([]string{"flower", "flow", "flight"}))
	fmt.Println(longestCommonPrefix2([]string{"dog", "racecar", "car"}))
}

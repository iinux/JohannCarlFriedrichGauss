package main

import "fmt"

func main() {
	fmt.Println(missingNumber([]int{3, 0, 1}) == 2)
	fmt.Println(missingNumber([]int{9, 6, 4, 2, 3, 5, 7, 0, 1}) == 8)
}

func missingNumber(nums []int) int {
	m := make(map[int]bool)
	for _, v := range nums {
		m[v] = true
	}
	l := len(nums)
	for i := 0; i <= l; i++ {
		_, ok := m[i]
		if !ok {
			return i
		}
	}

	return -1
}

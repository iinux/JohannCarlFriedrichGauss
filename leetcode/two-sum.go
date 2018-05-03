package main

import (
	"fmt"
)

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
}

func twoSum(nums []int, target int) []int {
	l := len(nums)
	var i, j int
	out:
	for i = 0; i < l; i++ {
		for j = i + 1; j < l; j++ {
			sum := nums[i] + nums[j]
			if sum == target {
				break out
			}
		}
	}

	return []int{i, j}
}

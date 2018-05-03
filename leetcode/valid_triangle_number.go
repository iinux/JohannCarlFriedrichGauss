package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(triangleNumber([]int{2, 2, 3, 4}))
}

func triangleNumber(nums []int) int {
	sort.Ints(nums)
	num_len := len(nums)
	num := 0
	for i := 0; i < num_len; i++ {
		for j := i + 1; j < num_len; j++ {
			for k := j + 1; k < num_len; k++ {
				if nums[k] >= nums[j] + nums[i] {
					break
				}
				num++
			}
		}
	}
	return num
}

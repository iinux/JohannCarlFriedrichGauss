package main

import "fmt"

func main() {
	fmt.Println(maxProduct([]int{2, 3, -2, 4}))
	fmt.Println(maxProduct([]int{-2, 0, -1}))
}

func maxProduct(nums []int) int {
	var max int
	var c int
	var cInit bool
	l := len(nums)
	for i := 0; i < l; i++ {
		c = nums[i]
		if !cInit || c > max {
			max = c
			cInit = true
		}
		for j := i + 1; j < l; j++ {
			c *= nums[j]
			if !cInit || c > max {
				max = c
				cInit = true
			}
		}
	}
	return max
}


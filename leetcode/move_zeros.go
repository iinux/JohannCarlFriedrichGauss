package main

import "fmt"

func main() {
	nums := []int{0, 1, 0, 3, 12}
	moveZeroes(nums)
	fmt.Println(nums)
}

func moveZeroes(nums []int) {
	l := len(nums)
	p := 0

	for i := 0; i < l; i++ {
		if nums[i] != 0 {
			nums[p] = nums[i]
			p++
		}
	}
	for ; p < l; p++ {
		nums[p] = 0
	}
}

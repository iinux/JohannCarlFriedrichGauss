package main

import "fmt"

func main() {
	fmt.Println(permute([]int{1, 2, 3}))
}

var l int
var use []bool
var res [][]int
var gNums []int

func permute(nums []int) [][]int {
	res = [][]int{}
	use = []bool{}
	gNums = nums
	l = len(nums)

	for i := 0; i < l; i++ {
		use = append(use, false)
	}

	calc([]int{})

	return res
}

func calc(current []int) {
	if len(current) == l {
		res = append(res, current)
		return
	}
	for i := 0; i < l; i++ {
		if !use[i] {
			a := make([]int, len(current))
			copy(a, current)
			a = append(a, gNums[i])
			use[i] = true
			calc(a)
			use[i] = false
		}
	}
}

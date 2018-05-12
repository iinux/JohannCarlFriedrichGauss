package main

import "fmt"

func main() {
	fmt.Println(combinationSum3(3, 7))
	fmt.Println(combinationSum3(3, 9))
}

var res [][]int

func combinationSum3(k int, n int) [][]int {
	res = [][]int{}
	calc([]int{}, 1, k, n)
	return res
}

func calc(current []int, start int, k int, n int) {
	if k == 1 {
		if n >= start && n <= 9 {
			current = append(current, n)
			res = append(res, current)
			return
		}
	}
	for i := start; i <= 9; i++ {
		a := make([]int, len(current))
		copy(a, current)
		a = append(a, i)
		calc(a, i + 1, k - 1, n - i)
	}
}

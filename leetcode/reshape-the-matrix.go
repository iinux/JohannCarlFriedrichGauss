package main

import "fmt"

func main() {
	fmt.Println(matrixReshape([][]int{{1, 2}, {3, 4}}, 1, 4))
	fmt.Println(matrixReshape([][]int{{1, 2}, {3, 4}}, 2, 4))
}

func matrixReshape(nums [][]int, r int, c int) [][]int {
	or := len(nums)
	oc := len(nums[0])

	if or * oc != r * c {
		return nums
	}

	var res [][]int
	var row []int
	var count int

	for i := 0; i < or; i++ {
		for j := 0; j < oc; j++ {
			row = append(row, nums[i][j])
			count++
			if count == c {
				res = append(res, row)
				row = []int{}
				count = 0
			}
		}
	}
	return res
}

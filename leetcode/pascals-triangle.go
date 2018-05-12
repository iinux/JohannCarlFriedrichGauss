package main

import "fmt"

func main() {
	fmt.Println(generate(5))
}

func generate(numRows int) [][]int {
	var r [][]int
	var t []int

	for i := 1; i <= numRows; i++ {
		if i == 1 {
			t = []int{1}
		} else if i == 2 {
			t = []int{1, 1}
		} else {
			var tt []int
			tt = append(tt, 1)
			tll := len(t) - 1
			for j := 0; j < tll; j++ {
				tt = append(tt, t[j] + t[j + 1])
			}
			tt = append(tt, 1)
			t = tt
		}
		r = append(r, t)
	}

	return r
}

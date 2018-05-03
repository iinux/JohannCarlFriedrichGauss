package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(validSquare([]int{0, 1}, []int{1, 1}, []int{1, 0}, []int{0, 1}))
	fmt.Println(validSquare([]int{0, 3}, []int{1, 1}, []int{1, 0}, []int{0, 1}))
	fmt.Println(validSquare([]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}))
}

func validSquare(p1 []int, p2 []int, p3 []int, p4 []int) bool {
	var lines []float64
	var ps [][]int
	ps = append(ps, p1)
	ps = append(ps, p2)
	ps = append(ps, p3)
	ps = append(ps, p4)

	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			lines = append(lines, math.Sqrt(math.Pow(float64(ps[i][0]-ps[j][0]), 2)+math.Pow(float64(ps[i][1]-ps[j][1]), 2)))
		}
	}

	var low, high float64
	var low_time, high_time int
	low = lines[0]
	high = lines[0]
	for _, l := range lines {
		if l < low {
			low = l
		}
		if l > high {
			high = l
		}
	}
	for _, l := range lines {
		if l == low {
			low_time++
		}
		if l == high {
			high_time++
		}
	}
	fmt.Println(lines)
	return low_time == 4 && high_time == 2

}

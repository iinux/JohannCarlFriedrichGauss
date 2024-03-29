package main

import (
	"fmt"
	"math"
)

func numSquares(n int) int {
	f := make([]int, n+1)
	for i := 1; i <= n; i++ {
		minn := math.MaxInt32
		for j := 1; j*j <= i; j++ {
			minn = min(minn, f[i-j*j])
		}
		f[i] = minn + 1
	}
	return f[n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println(3 == numSquares(12))
	fmt.Println(2 == numSquares(13))
	fmt.Println(numSquares(110))
	fmt.Println(numSquares(315))
}

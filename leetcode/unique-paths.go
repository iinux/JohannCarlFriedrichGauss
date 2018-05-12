package main

import (
	"fmt"
)

func main() {
	judge(3, 2, 3)
	judge(7, 3, 28)
	judge(51, 9, 0)
	judge(1, 10, 1)
}

func judge(m, n, a int) {
	myA := uniquePaths(m, n)
	fmt.Println(myA, a)
	fmt.Println(myA == a)
}

var cache map[string]int

func uniquePaths(m int, n int) int {
	cache = make(map[string]int)

	return calc(m, n)
}

func calc(m int, n int) int {
	var r int

	if m > n {
		m, n = n, m
	}

	key := fmt.Sprintf("%d,%d", m, n)
	r, ok := cache[key]
	if ok {
		return r
	}

	if m == 1 || n == 1 {
		r = 1
	} else if m == 2 {
		r = n
	} else if n == 2 {
		r = m
	} else {
		r = calc(m - 1, n) + calc(m, n - 1)
	}

	cache[key] = r

	return r
}


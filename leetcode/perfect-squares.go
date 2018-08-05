package main

import "fmt"

var rmap map[int]int

func numSquares(n int) int {
	if rmap == nil {
		rmap = make(map[int]int)
	}
	var ss []int
	for i := 1; true; i++ {
		if i*i > n {
			break
		} else {
			ss = append(ss, i*i)
		}
	}
	return numSquaresCalc(n, ss)
}

func numSquaresCalc(n int, c []int) int {
	if n <= 0 {
		return 0
	}
	var r int
	r, ok := rmap[n]
	if ok {
		//fmt.Println(rmap, n)
		//return r
	}

	var nc []int
	for _, v := range c {
		if v <= n {
			nc = append(nc, v)
		} else {
			break
		}
	}
	c = nc

	l := len(c)
	if l == 1 {
		return n
	}
	a := numSquaresCalc(n-c[l-1], c) + 1
	b := numSquaresCalc(n, c[0:l-1])
	if a > b {
		r = b
	} else {
		r = a
	}
	rmap[n] = r
	return r
}

func main() {
	fmt.Println(3 == numSquares(12))
	fmt.Println(2 == numSquares(13))
	fmt.Println(numSquares(110))
	fmt.Println(numSquares(315))
}

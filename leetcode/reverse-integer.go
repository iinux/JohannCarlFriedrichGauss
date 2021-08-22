package main

import "strconv"

func main()  {
	println(reverse(123))
	println(reverse(-123))
	println(reverse(120))
	println(reverse(0))
	println(reverse(1534236469))
	println(INT_MIN, INT_MAX)
}

const INT_MAX = int(^uint32(0) >> 1)
const INT_MIN = ^INT_MAX

func reverse(x int) int {
	lz := false
	if x < 0 {
		lz = true
		x = -x
	}

	y := []byte(strconv.Itoa(x))
	l := len(y)
	for i := 0 ; i < l; i++ {
		a := i
		b := l - i - 1
		if a > b {
			break
		}

		t := y[a]
		y[a] = y[b]
		y[b] = t
	}

	r, _ := strconv.Atoi(string(y))
	if lz {
		r = -r
	}
	if r > INT_MAX || r < INT_MIN {
		return 0
	}
	return r
}

package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(maximumSwap(2736))
	fmt.Println(maximumSwap(9973))
	fmt.Println(maximumSwap(1993))
}

func maximumSwap(num int) int {
	s := []rune(fmt.Sprintf("%d", num))
	l := len(s)
	for i := 0; i < l; i++ {
		var max rune
		maxi := -1
		for j := i + 1; j < l; j++ {
			if s[j] >= max {
				max = s[j]
				maxi = j
			}
		}
		if max > s[i] {
			s[i], s[maxi] = s[maxi], s[i]
			break
		}
	}

	v, _ := strconv.Atoi(string(s))

	return v
}

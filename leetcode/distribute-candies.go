package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(distributeCandies([]int{1, 1, 2, 2, 3, 3}) == 3)
	fmt.Println(distributeCandies([]int{1, 1, 2, 3}) == 2)
}

func distributeCandies(candies []int) int {
	sort.Ints(candies)
	num := len(candies) / 2
	p := -1
	c := 0
	for _, v := range candies {
		if c == num {
			break
		}
		if v != p {
			p = v
			c++
		}
	}

	return c
}

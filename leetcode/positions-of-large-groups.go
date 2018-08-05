package main

import "fmt"

func largeGroupPositions(S string) [][]int {
	r := make([][]int, 0)
	l := len(S)
	i := 0
	j := 1
	for i < l && j < l {
		if S[i] == S[j] {
			j++
		} else {
			if j - 1 - i + 1 >= 3 {
				r = append(r, []int{i, j - 1})
			}
			i = j
			j++
		}
	}

	if j - 1 - i + 1 >= 3 {
		r = append(r, []int{i, j - 1})
	}

	return r
}

func main()  {
	fmt.Println(largeGroupPositions("abbxxxxzzy"))
	fmt.Println(largeGroupPositions("abc"))
	fmt.Println(largeGroupPositions("aaa"))
}

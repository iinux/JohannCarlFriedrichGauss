package main

import "fmt"

func main() {
	s := make([]int, 2)
	mdSlice(s)
	// mdSlice2(&s)
	fmt.Println(s)
}

func mdSlice(s []int) {
	s = append(s, 1)
	s = append(s, 2)
	fmt.Println(s)
}

func mdSlice2(s *[]int) {
	*s = append(*s, 1)
	*s = append(*s, 2)
}
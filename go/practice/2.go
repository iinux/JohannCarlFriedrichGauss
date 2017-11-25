package main

import "fmt"

func main() {
	data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s := data[8:]
	s2 := data[:5]

	copy(s2, s) // dst:s2, src:s

	fmt.Println(s2)
	fmt.Println(data)
}

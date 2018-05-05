package main

import "fmt"

func main() {
	fmt.Println(isOneBitCharacter([]int{1, 0, 0}))
	fmt.Println(isOneBitCharacter([]int{1, 1, 1, 0}))
}

func isOneBitCharacter(bits []int) bool {
	l := len(bits)
	for i := 0; i < l; {
		if i == l - 1 && bits[i] == 0 {
			return true
		}
		if bits[i] == 0 {
			i++
		} else {
			i += 2
		}
	}
	return false
}

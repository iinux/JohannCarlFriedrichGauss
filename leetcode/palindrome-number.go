package main

import "strconv"

func main() {
	println(isPalindrome(121))
	println(isPalindrome(-121))
	println(isPalindrome(10))
}

func isPalindrome(x int) bool {
	y := strconv.Itoa(x)
	l := len(y)
	for i := 0; i < l/2; i++ {
		if y[i] != y[l-i-1] {
			return false
		}
	}

	return true
}

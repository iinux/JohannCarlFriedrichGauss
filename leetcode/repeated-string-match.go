package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(repeatedStringMatch("abcd", "cdabcda") == 3)
	fmt.Println(repeatedStringMatch("abcd", "cdabcdabcdab") == 4)
	fmt.Println(repeatedStringMatch("a", "aa") == 2)
	fmt.Println(repeatedStringMatch("abcd", "cdabcd") == 2)
	fmt.Println(repeatedStringMatch("aa", "a") == 1)
	fmt.Println(repeatedStringMatch("abababaaba", "aabaaba") == 2)
	fmt.Println(repeatedStringMatch("abcd", "abcdb") == -1)
	fmt.Println(repeatedStringMatch("abc", "aabcbabcc") == -1)
	fmt.Println(repeatedStringMatch("baa", "abaab") == 3)
}

func repeatedStringMatch(A string, B string) int {
	c := 1
	AA := A
	for true {
		if len(AA) > 2 * len(A) + len(B) {
			break
		}
		if strings.Index(AA, B) >= 0 {
			return c
		} else {
			AA += A
			c++
		}
	}

	return -1
}

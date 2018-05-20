package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(wordPattern("abba", "dog cat cat dog") == true)
	fmt.Println(wordPattern("abba", "dog cat cat fish") == false)
	fmt.Println(wordPattern("aaaa", "dog cat cat dog") == false)
	fmt.Println(wordPattern("abba", "dog dog dog dog") == false)
	fmt.Println(wordPattern("jquery", "jquery")==false)
}

func wordPattern(pattern string, str string) bool {
	strs := strings.Split(str, " ")
	wordMap := make(map[byte]string)
	l := len(pattern)
	sl := len(strs)
	if l != sl {
		return false
	}

	for i := 0; i < l; i++ {
		s, ok := wordMap[pattern[i]]
		if !ok {
			wordMap[pattern[i]] = strs[i]
		} else if s != strs[i] {
			return false
		}
	}

	strs = []string{}
	for _, v := range wordMap {
		strs = append(strs, v)
	}

	l = len(strs)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if strs[i] == strs[j] {
				return false
			}
		}
	}

	return true
}

package main

func main() {
	println(romanToInt("III"))
	println(romanToInt("IV"))
	println(romanToInt("IX"))
	println(romanToInt("LVIII"))
	println(romanToInt("MCMXCIV"))
}

const I = 1
const V = 5
const X = 10
const L = 50
const C = 100
const D = 500
const M = 1000

func romanToInt(s string) int {
	var s2 []int
	for _, e := range s {
		switch e {
		case 'I':
			s2 = append(s2, I)
		case 'V':
			s2 = append(s2, V)
		case 'X':
			s2 = append(s2, X)
		case 'L':
			s2 = append(s2, L)
		case 'C':
			s2 = append(s2, C)
		case 'D':
			s2 = append(s2, D)
		case 'M':
			s2 = append(s2, M)
		}
	}

	return calc(s2)
}

func calc(s2 []int) int {
	l := len(s2)
	if l == 1 {
		return s2[0]
	} else if l == 0 {
		return 0
	}

	maxIndex := -1
	max := 0
	for i, e := range s2 {
		if e > max {
			max = e
			maxIndex = i
		}
	}

	left := calc(s2[0:maxIndex])
	right := calc(s2[maxIndex+1:])
	return max - left + right
}

package main

func main() {
	println(isValid("()"))
	println(isValid("()[]{}"))
	println(isValid("(]"))
	println(isValid("([)]"))
	println(isValid("{[]}"))
	println(isValid("["))
}

func isValid(s string) bool {
	stack := []rune{}
	cmap := make(map[rune]rune)
	cmap[')'] = '('
	cmap[']'] = '['
	cmap['}'] = '{'
	l := 0
	for _, e := range s {
		switch e {
		case '(','[','{':
			stack = append(stack, e)
			l++
		case ')',']','}':
			if l < 1 {
				return false
			}
			if stack[l-1] != cmap[e] {
				return false
			} else {
				stack = stack[:l-1]
				l--
			}
		}
	}

	if l != 0 {
		return false
	}

	return true
}

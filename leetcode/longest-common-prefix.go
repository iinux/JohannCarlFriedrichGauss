package main

func main() {
	println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
	println(longestCommonPrefix([]string{"dog","racecar","car"}))
	println(longestCommonPrefix([]string{"cir","car"}))
}

func longestCommonPrefix(strs []string) string {
	l := len(strs)
	if l == 1 {
		return strs[0]
	} else if l == 0 {
		return ""
	}

	minLen := 200
	for i := 0; i < l; i++ {
		l2 := len(strs[i])
		if l2 < minLen {
			minLen = l2
		}
	}

	if minLen == 0 {
		return ""
	}

	n := 0

	for i := 0; i < minLen; i++ {
		eq := true
		for j := 1; j < l; j++ {
			if strs[0][i] != strs[j][i] {
				eq = false
				break
			}
		}

		if eq {
			n++
		} else {
			break
		}
	}

	return strs[0][:n]
}

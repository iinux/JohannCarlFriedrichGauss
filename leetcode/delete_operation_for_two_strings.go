package main

import (
	"fmt"
	//"log"
	"time"
)

func main() {
	judge(
		"sea",
		"eat", 2)
	judge(
		"park",
		"spake", 3)
	judge(
		"intention",
		"execution", 8)
	judge(
		"pvhvykrvntdywrhylaprgqmbzqitrhdmxboyw",
		"oelftlrthdmlwznwuritwrvdciho", 45)
	judge("zoologicoarchaeologist", "zoogeologist", 10)
	judge(
		"ytnxxvgngowdkosgcdiryvcozausbqatuxl",
		"pstzwvbktroytnzshnfrbnrhxaevezosydflydveggayzgr", 60)
	judge(
		"dmrnmzvbwibzkptlhuzmycuxcrxftwtgqaxyicjhswuwwvb",
		"wcdpnbjwkswzaajujigwuyiunpxexvgxxz", 57)
	judge(
		"ojmppcfzeupprdnuktgvtslgsjwxddvprjwtwiwattk",
		"vzjwkcwlerpt", 43)
	judge(
		"dinitrophenylhydrazine",
		"benzalphenylhydrazone", 13)
	judge(
		"mqmlmwfrmcavsxmnvhoovambqqukhcpigftohdvjjtirwiwv",
		"ukunnsxltoxgrrbkvhemvloyxmnfvnrvxprlhhvomoldhm", 68)
	judge(
		"abcdxabcde",
		"abcdeabcdx", 4)
	judge(
		"",
		"", 0)
	judge(
		"pneumonoultramicroscopicsilicovolcanoconiosis",
		"ultramicroscopic", 29)
	judge(
		"szwokpjlgqgogbgpjaczcmtjhzgldwinqfxbcxgghitcinmtdwnnpsmnmhfrhrgwncvcizaze",
		"spjlqggpzcgdxxtdwnrvca", 51)
	judge(
		"ebvivhpfxoptspwianmuhmkmbhxkqbrbgpfwpfcjixzhsjmtsgrzfshvkrvoxvjpmmsrojnpgzqdyofvicscopak",
		"vxoumkmxbpcixzhtrfhxmnzqyvisp", 59)
}

func judge(a string, b string, answer int) {
	startTime := time.Now()

	ya := minDistance(a, b)
	fmt.Println(ya, ya == answer, time.Now().Sub(startTime))
}

func minDistance(word1 string, word2 string) int {
	l1 := len(word1)
	l2 := len(word2)
	// return len(word1) + len(word2) - 2*lcs(word1, word2, len(word1), len(word2))
	var dp [][]int
	for i := 0; i <= l1; i++ {
		dp = append(dp, []int{})
		for j := 0; j <= l2; j++ {
			dp[i] = append(dp[i], 0)
		}
	}

	for i := 0; i <= l1; i++ {
		dp = append(dp, []int{0})
		for j := 0; j <= l2; j++ {
			if i == 0 || j == 0 {
				continue
			}
			if word1[i-1] == word2[j-1] {
				dp[i][j] = 1 + dp[i-1][j-1]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return l1 + l2 - 2 * dp[l1][l2]
}

func lcs(s string, s2 string, i int, i2 int) int {
	if i == 0 || i2 == 0 {
		return 0
	}

	if s[i-1] == s2[i2-1] {
		return 1 + lcs(s, s2, i-1, i2-1)
	} else {
		return max(lcs(s, s2, i, i2-1), lcs(s, s2, i-1, i2))
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

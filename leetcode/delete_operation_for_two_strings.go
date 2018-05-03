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
}

func judge(a string, b string, answer int) {
	startTime := time.Now()

	ya := minDistance(a, b)
	fmt.Println(ya, ya == answer, time.Now().Sub(startTime))
}

func minDistance(word1 string, word2 string) int {
	len1 := len(word1)
	len2 := len(word2)

	if len1 > 0 && len2 > 0 && word1[0] == word2[0] {
		return minDistance(word1[1:], word2[1:])
	}

	word2Map := make(map[int32][]int)
	for k, v := range word2 {
		word2Map[v] = append(word2Map[v], k)
	}

	var lines lineSlice
	for i := 0; i < len1; i++ {
		alpha := word1[i]
		list := word2Map[int32(alpha)]
		for _, v := range list {
			line := line{X: i, Y: v, Alpha: alpha}
			lines = append(lines, &line)
		}
	}

	len3 := calc(lines, 0, 0)

	return len1 + len2 - 2*len3
}

func calc(lines lineSlice, startX int, startY int) int {
	end := true
	var okLines lineSlice
	for _, line := range lines {
		if line.X >= startX && line.Y >= startY {
			okLines = append(okLines, line)
			end = false
		}
	}

	if end {
		return 0
	}

	max := 0
	for _, line := range okLines {
		v := calc(lines, line.X+1, line.Y+1) + 1
		if v > max {
			max = v
		}
	}
	return max
}

type line struct {
	X     int
	Y     int
	Alpha uint8
}

type lineSlice []*line

func minDistance_error(word1 string, word2 string) int {
	len1 := len(word1)
	len2 := len(word2)
	ret := 0
	for i := 0; i < len1; i++ {
		for j := 0; j < len2; j++ {
			l := comLen(word1[i:], word2[j:])
			if l > ret {
				ret = l
			}
		}
	}

	return len1 - ret + len2 - ret
}

func comLen(word1 string, word2 string) int {
	len1 := len(word1)
	len2 := len(word2)
	if len2 < len1 {
		len1 = len2
	}
	ret := 0
	for i := 0; i < len1; i++ {
		if word1[i] == word2[i] {
			ret++
		} else {
			break
		}
	}
	return ret
}

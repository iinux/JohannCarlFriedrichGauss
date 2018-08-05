package main

import (
	"fmt"
	//"log"
	"time"
	"sort"
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
	if len(word1) > len(word2) {
		word1, word2 = word2, word1
	}
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

	sort.Sort(lines)
	count = 0
	fmt.Println(len(lines))
	// len3 := calc(lines, 0, 0)
	len3 := calc2(lines, 0)
	fmt.Println(count)

	return len1 + len2 - 2*len3
}

var count int

func calc2(lines lineSlice, start int) int {
	var fi int
	end := true
	if start >= len(lines) {
		return 0
	}
	for i, line := range lines {
		if line.X >= lines[start].X && line.Y >= lines[start].Y {
			end = false
			fi = i
			break
		}
	}
	if end {
		return 0
	}
	a := calc2(lines, fi+1) + 1
	b := calc2(lines, start+1)
	if a >= b {
		return a
	} else {
		return b
	}
}

func calc(lines lineSlice, startX int, startY int) int {
	count++
	end := true
	//var okLines lineSlice
	minxi := -1
	minyi := -1
	for i, line := range lines {
		if line.X >= startX && line.Y >= startY {
			if minxi == -1 || line.X < lines[minxi].X || line.X == lines[minxi].X && line.Y < lines[minxi].Y {
				minxi = i
			}
			if minyi == -1 || line.Y < lines[minyi].Y || line.Y == lines[minyi].Y && line.X < lines[minyi].X {
				minyi = i
			}
			end = false
		}
	}
	/*
	for _, line := range lines {
		if line.X >= startX && line.Y >= startY {
			if line.X <= lines[minxi].X || line.Y <= lines[minxi].Y || line.X <= lines[minyi].X || line.Y <= lines[minyi].Y{
				okLines = append(okLines, line)
			}
			end = false
		}
	}
	*/

	if end {
		return 0
	}

	/*
	max := 0
	for _, line := range okLines {
		v := calc(lines, line.X+1, line.Y+1) + 1
		if v > max {
			max = v
		}
	}
	return max
	*/
	if minxi == minyi {
		return calc(lines, lines[minxi].X+1, lines[minxi].Y+1) + 1
	} else {
		var a []int
		a = append(a, calc(lines, lines[minxi].X+1, lines[minxi].Y+1)+1)
		a = append(a, calc(lines, lines[minyi].X+1, lines[minyi].Y+1)+1)
		a = append(a, calc(lines, lines[minxi].X+1, lines[minyi].Y+1))
		//a = append(a, calc(lines, startX, startY+1))
		max := 0
		for _, v := range a {
			if v > max {
				max = v
			}
		}
		return max
	}
}

type line struct {
	X     int
	Y     int
	Alpha uint8
}

type lineSlice []*line

func (c lineSlice) Len() int {
	return len(c)
}
func (c lineSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c lineSlice) Less(i, j int) bool {
	if c[i].X == c[j].X {
		return c[i].Y < c[j].Y
	} else {
		return c[i].X < c[j].X
	}
}

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

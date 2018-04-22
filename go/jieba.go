package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/wangbin/jiebago"
)

var seg jiebago.Segmenter

func init() {
	seg.LoadDictionary("G:/gopath/src/github.com/wangbin/jiebago/dict.txt")
}

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PairList) Len() int {
	return len(p)
}
func (p PairList) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

func main() {
	var err error
	var n int
	var isChinese bool
	freqMap := make(map[string]int)
	buf := make([]byte, 1000000)

	f, err := os.Open("a.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for true {
		n, err = f.Read(buf)
		if n < 1 {
			break
		}

		for word := range seg.Cut(string(buf[0:n]), false) {
			isChinese = true
			wordRune := []rune(word)

			if len(wordRune) < 2 {
				continue
			}

			for _, r := range wordRune {
				if r < 128 {
					isChinese = false
				}
			}
			if isChinese {
				freqMap[word]++
			}
		}
	}

	fmt.Println(sortMapByValue(freqMap))

	/*
	fmt.Print("【全模式】：")
	print(seg.CutAll("我来到北京清华大学"))

	fmt.Print("【精确模式】：")
	print(seg.Cut("我来到北京清华大学", false))

	fmt.Print("【新词识别】：")
	print(seg.Cut("他来到了网易杭研大厦", true))

	fmt.Print("【搜索引擎模式】：")
	print(seg.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))
	*/
}

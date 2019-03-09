package main

import (
	"fmt"
)

func main()  {
	judge([]int{5,3,4,5}, true)
	judge([]int{3,2,10,4}, true) // 局部最优不能保证全局最优，所以此用例没有通过，
}

func judge(params []int, a bool) {
	r := stoneGame(params)
	fmt.Println(r, a, r== a)
}

func stoneGame(piles []int) bool {
	alex := 0
	lee := 0
	p := true // true is alex false is lee
	l := 0
	num := 0
	for true {
		l = len(piles)
		if l <= 0 {
			break
		}
		fmt.Println(piles)
		if piles[0] > piles[l-1] {
			num = piles[0]
			piles = piles[1:]
		} else {
			num = piles[l-1]
			piles = piles[0:l-1]
		}
		fmt.Println(piles)
		if p {
			alex += num
		} else {
			lee += num
		}
		p = !p
	}

	return alex > lee
	// return true
	// 直接返回true竟然所有用例都通过，汗
}

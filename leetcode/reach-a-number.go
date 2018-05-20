package main

import "fmt"

func main() {
	judge(3, 2)
	judge(2, 3)
	judge(4, 3)
	judge(5, 5)
	judge(9, 5)
}

func judge(q int, a int) {
	ya := reachNumber(q)
	fmt.Printf("q:%d a:%d ya:%d r:%v\n", q, a, ya, a == ya)
}

func reachNumber(target int) int {
	if target < 0 {
		target = -target
	}
	n := 0
	pos := 0
	for true {
		if pos == target {
			return n
		} else if pos < target {
			n++
			pos += n
		} else {
			if (pos - target) % 2 == 0 {
				break
			} else {
				n++
				pos += n
			}
		}
	}

	return n

}

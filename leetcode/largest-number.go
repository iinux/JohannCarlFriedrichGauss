package main

import "fmt"

func main() {
	fmt.Println(largestNumber([]int{3, 30, 34, 5, 9}))
	fmt.Println(largestNumber([]int{121, 12}))
	fmt.Println(largestNumber([]int{1, 1, 1}))
	fmt.Println(largestNumber([]int{3, 43, 48, 94, 85, 33, 64, 32, 63, 66}))
	fmt.Println(largestNumber([]int{0, 0}))
}

func largestNumber(nums []int) string {
	l := len(nums)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if gt(nums[j], nums[i]) {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}
	s := ""
	for _, num := range nums {
		s += fmt.Sprintf("%d", num)
	}

	sr := []rune(s)
	srl := len(sr)
	for i := 0; i < srl; i++ {
		if sr[i] == '0' && i != srl - 1 {
			continue
		} else {
			sr = sr[i:]
			break
		}
	}

	return string(sr)
}

func gt(a int, b int) bool {
	ar := []rune(fmt.Sprintf("%d%d", a, b))
	br := []rune(fmt.Sprintf("%d%d", b, a))

	al := len(ar)
	for i := 0; i < al; i++ {
		if ar[i] > br[i] {
			return true
		} else if ar[i] < br[i] {
			return false
		}
	}

	return false
}


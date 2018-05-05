package main

import "fmt"

func main() {
	fmt.Println(largestNumber([]int{3, 30, 34, 5, 9}))
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
	return s
}

func gt(a int, b int) bool {
	ar := []rune(fmt.Sprintf("%d", a))
	br := []rune(fmt.Sprintf("%d", b))

	al := len(ar)
	bl := len(br)

	i := 0
	for true {
		if i >= al && i >= bl {
			break
		}
		if i >= al {
			if br[i - al] == br [i] {
			} else if br[i - al] < br [i] {
				return false
			} else {
				return true
			}
		} else if i >= bl {
			if ar[i - bl] == ar [i] {
			} else if ar[i - bl] < ar [i] {
				return true
			} else {
				return false
			}
		} else {
			if ar[i] == br[i] {
			} else if ar[i] < br[i] {
				return false
			} else {
				return true
			}
		}
		i++
	}

	return false
}


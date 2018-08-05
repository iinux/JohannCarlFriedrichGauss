package main

import "fmt"

func checkPossibility(nums []int) bool {
	l := len(nums)
	notNum := 0
	for i := 0; i < l - 1; i++ {
		if notNum > 1 {
			return false
		}
		if nums[i] > nums[i + 1] {
			notNum++
			if i + 2 <= l - 1 && i - 1 >= 0 && nums[i + 2] < nums[i] && nums[i + 1] < nums[i - 1] {
				return false
			}
		}
	}
	if notNum > 1 {
		return false
	}
	return true
}

func main() {
	fmt.Println(checkPossibility([]int{2, 3, 3, 2, 4}) == true)
	fmt.Println(checkPossibility([]int{4, 2, 3}) == true)
	fmt.Println(checkPossibility([]int{3, 4, 2, 3}) == false)
	fmt.Println(checkPossibility([]int{4, 2, 1}) == false)
}



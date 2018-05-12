package main

import "fmt"

func main() {
	fmt.Println(searchInsert([]int{1, 3, 5, 6}, 5) == 2)
	fmt.Println(searchInsert([]int{1, 3, 5, 6}, 2) == 1)
	fmt.Println(searchInsert([]int{1, 3, 5, 6}, 7) == 4)
	fmt.Println(searchInsert([]int{1, 3, 5, 6}, 0) == 0)

}

func searchInsert(nums []int, target int) int {
	for k, v := range nums {
		if target < v {
			return k
		} else if v == target {
			return k
		}
	}

	return len(nums)
}

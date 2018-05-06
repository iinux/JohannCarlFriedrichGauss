package main

import "fmt"

func main() {
	fmt.Println(findMaxAverage([]int{1, 12, -5, -6, 50, 3}, 4))

}

func findMaxAverage(nums []int, k int) float64 {
	l := len(nums)
	var max int
	maxInit := false
	for i := 0; i <= l - k; i++ {
		sum := 0
		for j := 0; j < k; j++ {
			sum += nums[i + j]
		}
		if !maxInit || sum > max {
			max = sum
			maxInit = true
		}
	}

	return float64(max) / float64(k)
}

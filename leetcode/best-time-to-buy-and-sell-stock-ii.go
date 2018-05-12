package main

import "fmt"

func main() {
	fmt.Println(maxProfit([]int{7, 1, 5, 3, 6, 4}) == 7)
	fmt.Println(maxProfit([]int{1, 2, 3, 4, 5}) == 4)
	fmt.Println(maxProfit([]int{7, 6, 4, 3, 1}) == 0)
	fmt.Println(maxProfit([]int{2, 1, 2, 0, 1}) == 2)
}

func maxProfit(prices []int) int {
	sum := 0
	stock := -1
	l := len(prices)
	for i := 0; i < l; i++ {
		if i == l - 1 {
			if stock >= 0 && prices[i] > stock {
				sum += prices[i] - stock
				stock = -1
			}
			break
		}
		if prices[i] > prices[i + 1] {
			if stock >= 0 {
				sum += prices[i] - stock
				stock = -1
			}
		} else if prices[i] == prices[i + 1] {
			continue
		} else {
			if stock == -1 {
				stock = prices[i]
			}
		}
	}

	return sum
}

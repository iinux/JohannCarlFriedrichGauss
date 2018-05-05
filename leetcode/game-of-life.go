package main

import "fmt"

func main() {
	var board [][]int
	board = append(board, []int{1, 1, 1})
	board = append(board, []int{1, 1, 1})
	board = append(board, []int{1, 1, 1})
	printBoard(board)
	gameOfLife(board)
	printBoard(board)
}

func printBoard(board [][]int)  {
	for k := range board {
		fmt.Println(board[k])
	}
	fmt.Println("== end ===")
}

var m, n int
var actions []*action

type action struct {
	x int
	y int
	v int
}

func gameOfLife(board [][]int) {
	m = len(board)
	n = len(board[0])
	actions = []*action{}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			c := calc(board, i, j)
			if c < 2 && board[i][j] == 1 {
				actions = append(actions, &action{x:i, y:j, v:0})
			} else if c > 3 && board[i][j] == 1 {
				actions = append(actions, &action{x:i, y:j, v:0})
			} else if c == 3 && board[i][j] == 0 {
				actions = append(actions, &action{x:i, y:j, v:1})
			}
		}
	}

	for _, a := range actions {
		board[a.x][a.y] = a.v
	}
}

func calc(board [][]int, x, y int) int {
	var c int
	c += calc2(board, x - 1, y - 1)
	c += calc2(board, x - 1, y)
	c += calc2(board, x - 1, y + 1)
	c += calc2(board, x, y - 1)
	c += calc2(board, x, y + 1)
	c += calc2(board, x + 1, y - 1)
	c += calc2(board, x + 1, y)
	c += calc2(board, x + 1, y + 1)

	return c
}

func calc2(board [][]int, x, y int) int {
	if x < 0 || y < 0 || x >= m || y >= n {
		return 0
	} else {
		return board[x][y]
	}
}

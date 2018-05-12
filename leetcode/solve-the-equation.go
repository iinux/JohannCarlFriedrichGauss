package main

import (
	"fmt"
	"strconv"
)

func main() {
	judge("x+5-3+x=6+x-2", "x=2")
	judge("-x=-1", "x=1")
	judge("0x=0", "Infinite solutions")
}

func judge(q, a string) bool {
	r := solveEquation(q) == a
	fmt.Println(r)
	return r
}

type exp struct {
	op    byte
	num   int
	power int
	numInit bool
}

func solveEquation(equation string) string {
	left := []*exp{}
	right := []*exp{}
	l := len(equation)
	isLeft := true
	segFinish := false
	t := &exp{
		op:'+',
	}
	for i := 0; i < l; i++ {
		if equation[i] == '=' {
			isLeft = false
			t.op = '+'
			continue
		}

		if i == l - 1 ||
			equation[i + 1] == '+' ||
			equation[i + 1] == '-' ||
			equation[i + 1] == '=' {
			segFinish = true
		}

		if equation[i] == 'x' {
			t.power = 1
			if !t.numInit {
				t.num = 1
			}
		} else if equation[i] == '+' || equation[i] == '-' {
			t.op = equation[i]
		} else {
			tt, err := strconv.Atoi(string(equation[i]))
			if err != nil {
				fmt.Println(err)
			}
			t.num = t.num * 10 + tt
			t.numInit = true
		}
		if segFinish {
			if isLeft {
				left = append(left, t)
			} else {
				right = append(right, t)
			}
			t = &exp{}
			segFinish = false
		}
	}

	left2 := exp{power:1}
	left3 := exp{}
	for _, v := range left {
		if v.op == '-' {
			v.num = -v.num
		}
		if v.power == 1 {
			left2.num += v.num
		} else {
			left3.num += v.num
		}
	}

	right2 := exp{power:1}
	right3 := exp{}
	for _, v := range right {
		if v.op == '-' {
			v.num = -v.num
		}
		if v.power == 1 {
			right2.num += v.num
		} else {
			right3.num += v.num
		}
	}

	xSum := left2.num - right2.num
	numSum := right3.num - left3.num

	for _, v := range left {
		fmt.Println(v)
	}
	for _, v := range right {
		fmt.Println(v)
	}
	fmt.Println(left2)
	fmt.Println(left3)
	fmt.Println(right2)
	fmt.Println(right3)
	fmt.Println(xSum, numSum)
	fmt.Println("======")

	if xSum == 0 {
		if numSum == 0 {
			return "Infinite solutions"
		} else {
			return "No solution"
		}
	} else {
		return fmt.Sprintf("x=%d", numSum / xSum)
	}

	return ""
}

package main

import (
	"fmt"
)

type Request struct {
	data []int
	ret  chan int
}

func NewRequest(data ...int) *Request {
	return &Request{data, make(chan int, 1)}
}

func Process(req *Request) {
	x := 0
	for _, i := range req.data {
		x += i
	}

	req.ret <- x
}

func main() {
	req := NewRequest(10, 20, 30)
	Process(req)
	fmt.Println(<-req.ret)
}

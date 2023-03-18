package main

import (
"fmt"
"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println(s)
	}
}

func main() {
	fmt.Println("Go Version:", runtime.Version())
	go say("world") //开一个新的Goroutines执行
	say("hello") //当前Goroutines执行
}

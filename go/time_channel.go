package main

import (
	"fmt"
	"time"
)

func main()  {
	tick := time.Tick(1000 * time.Millisecond)
	for {
		<-tick
		fmt.Println("tick")
	}
}

func example() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
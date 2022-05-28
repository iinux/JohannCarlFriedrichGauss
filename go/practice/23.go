package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer fmt.Println("A.defer")

		func() {
			defer fmt.Println("B.defer")
			runtime.Goexit() // 终止当前 goroutine
			runtime.Version()
			fmt.Println("B") // 不会执行
		}()

		fmt.Println("A") // 不会执行
	}()
	fmt.Println("main end")

	wg.Wait()
}

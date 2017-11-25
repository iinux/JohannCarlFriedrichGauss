package main

import "fmt"

func main() {
	data := make(chan int)  // 数据交换队列
	exit := make(chan bool) // 退出通知

	go func() {
		for d := range data { // 从队列迭代接收数据，直到 close 。
			fmt.Println(d)
		}
		fmt.Println("recv over.")
		exit <- true // 发出退出通知。
	}()

	data <- 1 // 发送数据。
	data <- 2
	data <- 3
	close(data) // 关闭队列。

	fmt.Println("send over.")
	<-exit // 等待退出通知。
}

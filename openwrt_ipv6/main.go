package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	PORT = ":5800"
	FILE = "/tmp/openwrt_ipv6"
)

func main() {
	// 创建 UDP6 地址，监听所有 IPv6 接口的 5800 端口
	addr, err := net.ResolveUDPAddr("udp6", PORT)
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}

	// 创建 UDP 连接
	conn, err := net.ListenUDP("udp6", addr)
	if err != nil {
		log.Fatalf("Error listening on UDP: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Listening on UDP port %s (IPv6 only)\n", PORT)

	// 创建缓冲区
	buffer := make([]byte, 1024)

	for {
		// 读取 UDP 数据
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP message: %v", err)
			continue
		}

		// 将接收到的数据转换为字符串
		message := string(buffer[:n])
		
		// 打印接收到的消息信息
		fmt.Printf("Received %d bytes from %s\n", n, clientAddr)
		
		// 将数据写入文件（覆盖模式）
		err = os.WriteFile(FILE, []byte(message), 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			continue
		}
		
		fmt.Printf("Successfully wrote %d bytes to %s\n", n, FILE)
	}
}

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("Usage: go run main.go <port> <size> <count>")
		return
	}

	size, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid size:", err)
		return
	}

	count, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid count:", err)
		return
	}

	sizeBytes := make([]byte, 4)
	for i := range sizeBytes {
		sizeBytes[i] = byte((size >> (8 * (3 - i))) & 0xFF)
	}

	serverAddr := "localhost:" + args[0]

	for i := 0; i < count; i++ {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			os.Exit(1)
		}
		defer conn.Close()

		message := append([]byte{0x80, 0x01, 0x00, 0x01}, sizeBytes...)

		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error sending data:", err)
			os.Exit(1)
		}

		fmt.Println("Data sent successfully!")
	}
}

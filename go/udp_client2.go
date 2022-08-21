package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	flag.Parse()
	laddr := &net.UDPAddr{Port: 8021}
	raddr1, err := net.ResolveUDPAddr("udp", "hw.iinux.cn:8021")
	raddr2, err := net.ResolveUDPAddr("udp", "rack4.iinux.cn:8021")
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", laddr, raddr1)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}

	_, err = conn.Write([]byte("hi"))
	if err != nil {
		fmt.Println("write failed:", err)
		os.Exit(1)
	}
	data := make([]byte, 32)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}
	fmt.Println(string(data))

	_ = conn.Close()

	conn, err = net.DialUDP("udp", laddr, raddr2)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	_, err = conn.Write([]byte("hi"))
	if err != nil {
		fmt.Println("write failed:", err)
		os.Exit(1)
	}
	data = make([]byte, 32)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}
	fmt.Println(string(data))

	os.Exit(0)
}

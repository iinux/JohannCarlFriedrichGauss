package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "37", "port")

//go run ucp_client.go -host time.nist.gov
func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Can't close: ", err)
		}
	}(conn)
	_, err = conn.Write([]byte(""))
	if err != nil {
		fmt.Println("write failed:", err)
		os.Exit(1)
	}
	data := make([]byte, 4)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}
	t := binary.BigEndian.Uint32(data)
	fmt.Println(time.Unix(int64(t), 0).String())
	os.Exit(0)
}

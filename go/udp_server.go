package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "37", "port")

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing:", err)
		}
	}(conn)
	for {
		handleClient(conn)
	}
}
func handleClient(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}
	daytime := time.Now().Unix()
	fmt.Println(n, remoteAddr, string(data))
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(daytime))
	b1 := []byte(fmt.Sprintf("%d %s", daytime, remoteAddr))
	_, err = conn.WriteToUDP(b1, remoteAddr)
	if err != nil {
		fmt.Println("write to udp error:", err)
		return
	}
}

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"github.com/blevesearch/bleve/size"
	"encoding/json"
	"encoding/binary"
	"time"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "3737", "port")
var code = flag.String("code", "1234", "auth code")

func main() {
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *host + ":" + *port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}

	_, err = conn.Write([]byte(*code))
	if err != nil {
		fmt.Println("failed:", err)
		os.Exit(1)
	}

	data := make([]byte, size.SizeOfUint16)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}

	l := binary.BigEndian.Uint16(data)
	data = make([]byte, l)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}

	var peerAddr net.UDPAddr
	err = json.Unmarshal(data, &peerAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(peerAddr)

	go func() {
		data = make([]byte, 5)
		fmt.Println("before read")
		dataLen, err := conn.Read(data)
		fmt.Println("after read")
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err)
			return
		}
		fmt.Println(dataLen, string(data))
	}()

	/*
	for true {
		fmt.Println("before write")
		conn.WriteToUDP([]byte("hello"), &peerAddr)
		fmt.Println("after write")
		if err != nil {
			fmt.Println("failed to write UDP msg because of ", err)
		}
		time.Sleep(time.Minute)
	}
	*/

	peerConn, err := net.DialUDP("udp", nil, &peerAddr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}

	go func() {
		for true {
			_, err = peerConn.Write([]byte(*code))
			if err != nil {
				fmt.Println("failed:", err)
				os.Exit(1)
			}
			time.Sleep(time.Minute)
		}
	}()

	go func() {
		for true {
			data = make([]byte, len(*code))
			fmt.Println("before peer read")
			_, err = peerConn.Read(data)
			fmt.Println("after peer read")
			if err != nil {
				fmt.Println("failed to read UDP msg because of ", err)
				os.Exit(1)
			}
		}
	}()

	for true {
		time.Sleep(time.Minute)
	}
}


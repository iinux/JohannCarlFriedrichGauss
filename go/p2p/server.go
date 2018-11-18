package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"encoding/json"
	"encoding/binary"
	"github.com/blevesearch/bleve/size"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "3737", "port")

type pair struct {
	UdpAddr *net.UDPAddr
	Ch      chan *net.UDPAddr
}

var pairMap map[string]pair

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", *host + ":" + *port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()
	pairMap = make(map[string]pair)
	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	fourDigits := make([]byte, 4)
	_, remoteAddr, err := conn.ReadFromUDP(fourDigits)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}

	go func() {
		key := string(fourDigits)
		p, ok := pairMap[key]
		var peerAddr *net.UDPAddr
		if ok {
			fmt.Println("ok")
			p.Ch <- remoteAddr
			peerAddr = p.UdpAddr
			delete(pairMap, key)
		} else {
			fmt.Println("not ok")
			p = pair{
				UdpAddr:remoteAddr,
				Ch:make(chan *net.UDPAddr),
			}
			pairMap[key] = p
			peerAddr = <-p.Ch
		}
		fmt.Println(peerAddr)
		data, err := json.Marshal(peerAddr)
		if err != nil {
			fmt.Println(err)
			return
		}

		l := len(data)
		b := make([]byte, size.SizeOfUint16)
		binary.BigEndian.PutUint16(b, uint16(l))

		conn.WriteToUDP(b, remoteAddr)
		conn.WriteToUDP(data, remoteAddr)
	}()
}


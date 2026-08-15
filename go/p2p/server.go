package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "3737", "port")

type pair struct {
	UdpAddr *net.UDPAddr
}

var pairMap map[string]pair

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
	defer conn.Close()
	pairMap = make(map[string]pair)
	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	code := make([]byte, 256)
	n, remoteAddr, err := conn.ReadFromUDP(code)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}

	key := string(code[:n])
	p, ok := pairMap[key]
	if !ok {
		pairMap[key] = pair{UdpAddr: remoteAddr}
		fmt.Printf("waiting for peer: code=%q address=%s\n", key, remoteAddr)
		return
	}

	delete(pairMap, key)
	fmt.Printf("paired code=%q: %s <-> %s\n", key, p.UdpAddr, remoteAddr)
	sendPeerAddress(conn, remoteAddr, p.UdpAddr)
	sendPeerAddress(conn, p.UdpAddr, remoteAddr)
}

func sendPeerAddress(conn *net.UDPConn, destination, peerAddr *net.UDPAddr) {
	data, err := json.Marshal(peerAddr)
	if err != nil {
		fmt.Println("failed to encode peer address:", err)
		return
	}
	if _, err := conn.WriteToUDP(data, destination); err != nil {
		fmt.Println("failed to write UDP msg because of", err)
	}
}

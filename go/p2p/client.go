package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "3737", "port")
var code = flag.String("code", "1234", "auth code")

func main() {
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	// Use one unconnected UDP socket for both the rendezvous server and the
	// peer. This keeps the same local port (and therefore the same NAT mapping)
	// for every packet.
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("Can't listen: ", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("local UDP:", conn.LocalAddr())
	fmt.Println("rendezvous server:", addr)

	_, err = conn.WriteToUDP([]byte(*code), addr)
	if err != nil {
		fmt.Println("failed:", err)
		os.Exit(1)
	}

	data := make([]byte, 64*1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err)
		os.Exit(1)
	}
	if !sameUDPAddr(remoteAddr, addr) {
		fmt.Println("received unexpected UDP message from", remoteAddr)
		os.Exit(1)
	}

	var peerAddr net.UDPAddr
	err = json.Unmarshal(data[:n], &peerAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("peer public address:", &peerAddr)

	peerReady := make(chan struct{})

	go func() {
		const punchInterval = 500 * time.Millisecond
		const punchTimeout = 30 * time.Second
		const keepaliveInterval = 20 * time.Second

		send := func() bool {
			if _, err := conn.WriteToUDP([]byte(*code), &peerAddr); err != nil {
				fmt.Println("failed to write peer UDP msg:", err)
				return false
			}
			return true
		}

		punchTicker := time.NewTicker(punchInterval)
		punchTimer := time.NewTimer(punchTimeout)
		defer punchTicker.Stop()
		defer punchTimer.Stop()

		fmt.Println("starting UDP hole punching")
		if !send() {
			return
		}

	punching:
		for {
			select {
			case <-peerReady:
				fmt.Println("UDP hole punching succeeded")
				break punching
			case <-punchTimer.C:
				fmt.Println("UDP hole punching timed out; continuing with keepalives")
				break punching
			case <-punchTicker.C:
				if !send() {
					return
				}
			}
		}

		keepaliveTicker := time.NewTicker(keepaliveInterval)
		defer keepaliveTicker.Stop()

		for range keepaliveTicker.C {
			if !send() {
				return
			}
		}
	}()

	peerSeen := false
	for {
		data = make([]byte, 64*1024)
		dataLen, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("failed to read UDP msg because of ", err)
			return
		}
		if !sameUDPAddr(remoteAddr, &peerAddr) {
			fmt.Printf("ignored %d bytes from unexpected address %s\n", dataLen, remoteAddr)
			continue
		}
		if !peerSeen {
			peerSeen = true
			close(peerReady)
		}
		fmt.Printf("received %d bytes from peer %s: %q\n", dataLen, remoteAddr, data[:dataLen])
	}
}

func sameUDPAddr(a, b *net.UDPAddr) bool {
	return a != nil && b != nil && a.Port == b.Port && a.IP.Equal(b.IP)
}

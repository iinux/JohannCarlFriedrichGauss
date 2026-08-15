package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

const matchProtocol = protocol.ID("/iinux/p2p/match/1.0.0")

type matchRequest struct {
	Code string `json:"code"`
}

type matchResponse struct {
	PeerID string `json:"peer_id,omitempty"`
	Error  string `json:"error,omitempty"`
}

type waiter struct {
	peerID peer.ID
	stream network.Stream
}

type matcher struct {
	mu      sync.Mutex
	waiting map[string]waiter
}

func main() {
	listen := flag.String("listen", "/ip4/0.0.0.0/tcp/4001", "relay listen multiaddr")
	announce := flag.String("announce", "", "public multiaddr, for example /dns4/hw.iinux.cn/tcp/4001")
	flag.Parse()

	listenAddr, err := ma.NewMultiaddr(*listen)
	if err != nil {
		log.Fatal("invalid listen address: ", err)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrs(listenAddr),
		libp2p.EnableRelayService(),
		libp2p.ForceReachabilityPublic(),
	}
	if *announce != "" {
		announceAddr, err := ma.NewMultiaddr(*announce)
		if err != nil {
			log.Fatal("invalid announce address: ", err)
		}
		opts = append(opts, libp2p.AddrsFactory(func([]ma.Multiaddr) []ma.Multiaddr {
			return []ma.Multiaddr{announceAddr}
		}))
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	m := &matcher{waiting: make(map[string]waiter)}
	h.SetStreamHandler(matchProtocol, m.handle)

	fmt.Println("relay peer ID:", h.ID())
	for _, addr := range h.Addrs() {
		fmt.Printf("client relay address: %s/p2p/%s\n", addr, h.ID())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

func (m *matcher) handle(stream network.Stream) {
	requester := stream.Conn().RemotePeer()
	var req matchRequest
	if err := json.NewDecoder(bufio.NewReader(stream)).Decode(&req); err != nil {
		writeResponse(stream, matchResponse{Error: "invalid request"})
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	if req.Code == "" {
		writeResponse(stream, matchResponse{Error: "code cannot be empty"})
		return
	}

	m.mu.Lock()
	previous, found := m.waiting[req.Code]
	if !found {
		m.waiting[req.Code] = waiter{peerID: requester, stream: stream}
		m.mu.Unlock()
		fmt.Printf("waiting: code=%q peer=%s\n", req.Code, requester)
		return
	}
	delete(m.waiting, req.Code)
	m.mu.Unlock()

	if previous.peerID == requester {
		writeResponse(stream, matchResponse{Error: "same peer cannot match itself"})
		writeResponse(previous.stream, matchResponse{Error: "duplicate peer"})
		return
	}

	fmt.Printf("paired: code=%q %s <-> %s\n", req.Code, previous.peerID, requester)
	writeResponse(previous.stream, matchResponse{PeerID: requester.String()})
	writeResponse(stream, matchResponse{PeerID: previous.peerID.String()})
}

func writeResponse(stream network.Stream, response matchResponse) {
	defer stream.Close()
	if err := json.NewEncoder(stream).Encode(response); err != nil {
		_ = stream.Reset()
	}
}

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	relayclient "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	matchProtocol = protocol.ID("/iinux/p2p/match/1.0.0")
	chatProtocol  = protocol.ID("/iinux/p2p/chat/1.0.0")
)

type matchRequest struct {
	Code string `json:"code"`
}

type matchResponse struct {
	PeerID string `json:"peer_id,omitempty"`
	Error  string `json:"error,omitempty"`
}

func main() {
	relayFlag := flag.String("relay", "", "relay multiaddr including /p2p/<peer-id>")
	code := flag.String("code", "1234", "pairing code")
	message := flag.String("message", "hello from libp2p", "message sent to the peer")
	flag.Parse()
	if *relayFlag == "" {
		log.Fatal("-relay is required; copy 'client relay address' from the server")
	}

	relayAddr, err := ma.NewMultiaddr(*relayFlag)
	if err != nil {
		log.Fatal("invalid relay address: ", err)
	}
	relayInfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		log.Fatal("relay address must end with /p2p/<peer-id>: ", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	h, err := libp2p.New(
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.ForceReachabilityPrivate(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()
	h.SetStreamHandler(chatProtocol, receiveChat)
	fmt.Println("local peer ID:", h.ID())

	if err := h.Connect(ctx, *relayInfo); err != nil {
		log.Fatal("connect relay: ", err)
	}
	reservation, err := relayclient.Reserve(ctx, h, *relayInfo)
	if err != nil {
		log.Fatal("reserve relay slot: ", err)
	}
	fmt.Println("relay reservation expires:", reservation.Expiration.Format(time.RFC3339))

	peerID, err := findPeer(ctx, h, relayInfo.ID, *code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("matched peer:", peerID)

	circuitSuffix, _ := ma.NewMultiaddr("/p2p-circuit/p2p/" + peerID.String())
	peerRelayAddr := relayAddr.Encapsulate(circuitSuffix)
	peerInfo, err := peer.AddrInfoFromP2pAddr(peerRelayAddr)
	if err != nil {
		log.Fatal("build peer relay address: ", err)
	}
	if err := h.Connect(ctx, *peerInfo); err != nil {
		log.Fatal("connect peer through relay: ", err)
	}
	fmt.Println("connected through relay; libp2p will attempt a direct upgrade")

	go sendChat(ctx, h, peerID, *message)
	go reportPath(ctx, h, peerID)
	<-ctx.Done()
}

func findPeer(ctx context.Context, h interface {
	NewStream(context.Context, peer.ID, ...protocol.ID) (network.Stream, error)
}, relayID peer.ID, code string) (peer.ID, error) {
	stream, err := h.NewStream(ctx, relayID, matchProtocol)
	if err != nil {
		return "", fmt.Errorf("open matching stream: %w", err)
	}
	defer stream.Close()
	if err := json.NewEncoder(stream).Encode(matchRequest{Code: code}); err != nil {
		return "", fmt.Errorf("send matching code: %w", err)
	}

	var response matchResponse
	if err := json.NewDecoder(bufio.NewReader(stream)).Decode(&response); err != nil {
		return "", fmt.Errorf("wait for peer: %w", err)
	}
	if response.Error != "" {
		return "", fmt.Errorf("matching failed: %s", response.Error)
	}
	return peer.Decode(response.PeerID)
}

func sendChat(ctx context.Context, h interface {
	NewStream(context.Context, peer.ID, ...protocol.ID) (network.Stream, error)
}, peerID peer.ID, message string) {
	streamCtx := network.WithAllowLimitedConn(ctx, "chat over relay")
	stream, err := h.NewStream(streamCtx, peerID, chatProtocol)
	if err != nil {
		log.Println("open chat stream:", err)
		return
	}
	defer stream.Close()
	if _, err := fmt.Fprintln(stream, message); err != nil {
		log.Println("send chat message:", err)
	}
}

func receiveChat(stream network.Stream) {
	defer stream.Close()
	data, err := io.ReadAll(io.LimitReader(stream, 64*1024))
	if err != nil {
		log.Println("read chat message:", err)
		return
	}
	fmt.Printf("message from %s: %s\n", stream.Conn().RemotePeer(), strings.TrimSpace(string(data)))
}

func reportPath(ctx context.Context, h interface {
	Network() network.Network
}, peerID peer.ID) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	last := ""
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			path := "disconnected"
			for _, conn := range h.Network().ConnsToPeer(peerID) {
				candidate := "direct: " + conn.RemoteMultiaddr().String()
				if strings.Contains(conn.RemoteMultiaddr().String(), "/p2p-circuit") {
					candidate = "relay: " + conn.RemoteMultiaddr().String()
				}
				if path == "disconnected" || strings.HasPrefix(candidate, "direct:") {
					path = candidate
				}
			}
			if path != last {
				fmt.Println("connection path:", path)
				last = path
			}
		}
	}
}

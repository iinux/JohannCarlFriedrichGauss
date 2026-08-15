package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
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
	proxyProtocol = protocol.ID("/iinux/p2p/tcp-proxy/1.0.0")
)

type matchRequest struct {
	Code string `json:"code"`
	Role string `json:"role"`
}

type matchResponse struct {
	PeerID string `json:"peer_id,omitempty"`
	Error  string `json:"error,omitempty"`
}

func main() {
	relayFlag := flag.String("relay", "", "relay multiaddr including /p2p/<peer-id>")
	code := flag.String("code", "1234", "pairing code")
	targetIP := flag.String("target-ip", "", "client1 target host or IP")
	targetPort := flag.Int("target-port", 0, "client1 target TCP port")
	listenIP := flag.String("listen-ip", "127.0.0.1", "client2 local listen IP")
	listenPort := flag.Int("listen-port", 0, "client2 local TCP listen port")
	flag.Parse()
	if *relayFlag == "" {
		log.Fatal("-relay is required; copy 'client relay address' from the server")
	}
	role, err := clientRole(*targetIP, *targetPort, *listenPort)
	if err != nil {
		log.Fatal(err)
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
	fmt.Printf("local peer ID: %s (role=%s)\n", h.ID(), role)

	if err := h.Connect(ctx, *relayInfo); err != nil {
		log.Fatal("connect relay: ", err)
	}
	reservation, err := relayclient.Reserve(ctx, h, *relayInfo)
	if err != nil {
		log.Fatal("reserve relay slot: ", err)
	}
	fmt.Println("relay reservation expires:", reservation.Expiration.Format(time.RFC3339))

	peerID, err := findPeer(ctx, h, relayInfo.ID, *code, role)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("matched peer:", peerID)
	if role == "provider" {
		target := net.JoinHostPort(*targetIP, strconv.Itoa(*targetPort))
		h.SetStreamHandler(proxyProtocol, func(stream network.Stream) {
			if stream.Conn().RemotePeer() != peerID {
				_ = stream.Reset()
				return
			}
			go forwardToTarget(stream, target)
		})
	}

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

	go reportPath(ctx, h, peerID)
	if role == "visitor" {
		listenAddress := net.JoinHostPort(*listenIP, strconv.Itoa(*listenPort))
		if err := serveLocalTCP(ctx, h, peerID, listenAddress); err != nil {
			log.Fatal(err)
		}
		return
	}
	fmt.Printf("provider ready: forwarding peer streams to %s\n", net.JoinHostPort(*targetIP, strconv.Itoa(*targetPort)))
	<-ctx.Done()
}

func clientRole(targetIP string, targetPort, listenPort int) (string, error) {
	hasTarget := targetIP != "" || targetPort != 0
	hasListener := listenPort != 0
	if hasTarget == hasListener {
		return "", fmt.Errorf("choose exactly one mode: client1 uses -target-ip and -target-port; client2 uses -listen-port")
	}
	if hasTarget {
		if targetIP == "" || targetPort < 1 || targetPort > 65535 {
			return "", fmt.Errorf("client1 requires a valid -target-ip and -target-port")
		}
		return "provider", nil
	}
	if listenPort < 1 || listenPort > 65535 {
		return "", fmt.Errorf("invalid -listen-port")
	}
	return "visitor", nil
}

func findPeer(ctx context.Context, h interface {
	NewStream(context.Context, peer.ID, ...protocol.ID) (network.Stream, error)
}, relayID peer.ID, code, role string) (peer.ID, error) {
	stream, err := h.NewStream(ctx, relayID, matchProtocol)
	if err != nil {
		return "", fmt.Errorf("open matching stream: %w", err)
	}
	defer stream.Close()
	if err := json.NewEncoder(stream).Encode(matchRequest{Code: code, Role: role}); err != nil {
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

func serveLocalTCP(ctx context.Context, h interface {
	NewStream(context.Context, peer.ID, ...protocol.ID) (network.Stream, error)
}, peerID peer.ID, listenAddress string) error {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", listenAddress, err)
	}
	defer listener.Close()
	fmt.Println("client2 listening on:", listener.Addr())
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Println("accept local TCP connection:", err)
			continue
		}
		go forwardToPeer(ctx, h, peerID, conn.(*net.TCPConn))
	}
}

func forwardToPeer(ctx context.Context, h interface {
	NewStream(context.Context, peer.ID, ...protocol.ID) (network.Stream, error)
}, peerID peer.ID, local *net.TCPConn) {
	streamCtx := network.WithAllowLimitedConn(ctx, "TCP proxy over relay")
	stream, err := h.NewStream(streamCtx, peerID, proxyProtocol)
	if err != nil {
		log.Println("open proxy stream:", err)
		_ = local.Close()
		return
	}
	log.Printf("proxy opened: %s -> peer %s", local.RemoteAddr(), peerID)
	bridge(local, stream)
}

func forwardToTarget(stream network.Stream, target string) {
	targetConn, err := net.DialTimeout("tcp", target, 10*time.Second)
	if err != nil {
		log.Printf("connect target %s: %v", target, err)
		_ = stream.Reset()
		return
	}
	log.Printf("proxy opened: peer %s -> %s", stream.Conn().RemotePeer(), target)
	bridge(targetConn.(*net.TCPConn), stream)
}

func bridge(tcpConn *net.TCPConn, stream network.Stream) {
	defer tcpConn.Close()
	defer stream.Close()
	done := make(chan struct{}, 2)
	go func() {
		_, _ = io.Copy(stream, tcpConn)
		_ = stream.CloseWrite()
		done <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(tcpConn, stream)
		_ = tcpConn.CloseWrite()
		done <- struct{}{}
	}()
	<-done
	<-done
	log.Println("proxy closed")
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

// natpunch is a minimal UDP NAT hole-punching demo.
//
// Two peers, each behind (potentially different) NATs, run the same `punch`
// command. A third host with a reachable public address runs `server`. The
// signaling server learns each peer's public UDP endpoint (by receiving a
// probe datagram from each peer) and introduces the two peers to each other.
// The peers then simultaneously send UDP packets to each other, which causes
// each NAT to install a mapping for the other's address — once each side has
// both sent AND received a packet from the other, the two peers can talk
// directly peer-to-peer, with the server out of the loop.
//
// Build:
//	go build -o natpunch .
//
// Local demo (no real NAT involved, but exercises the full protocol):
//
//	./natpunch server -listen 127.0.0.1:9000 -udp 127.0.0.1:9001
//	./natpunch punch  -server 127.0.0.1:9000 -udp-server 127.0.0.1:9001 \
//	                  -session room1 -name alice -local 127.0.0.1:7000
//	./natpunch punch  -server 127.0.0.1:9000 -udp-server 127.0.0.1:9001 \
//	                  -session room1 -name bob   -local 127.0.0.1:7001
//
// In real use, run `server` on a host with a public IP, and the two `punch`
// clients on two different NAT'd networks pointing at it. The server only
// needs to be reachable for signaling; once the punch lands, peers talk
// directly to each other.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// =====================================================================
// Wire protocol
// =====================================================================

// Msg is the JSON-over-TCP signaling message exchanged with the server.
type Msg struct {
	Type    string `json:"type"`
	Session string `json:"session,omitempty"`
	Name    string `json:"name,omitempty"`
	Addr    string `json:"addr,omitempty"` // "host:port" of a peer's public UDP endpoint
	Text    string `json:"text,omitempty"`
}

const (
	msgRegister   = "register"
	msgRegistered = "registered"
	msgPeer       = "peer"
	msgError      = "error"
)

// Datagram format on the peer-to-peer UDP channel:
//   "PUNCH|<name>|<seq>" — keep-alive heartbeat (one is sent every 200ms)
//   "CHAT|<name>|<text>" — a chat message
//
// The first datagram received from the peer (regardless of format) is what
// confirms the NAT mappings are open in both directions.

// =====================================================================
// Entry point
// =====================================================================

func main() {
	log.SetFlags(log.Lmicroseconds)
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "server":
		runServer(os.Args[2:])
	case "punch":
		runPunch(os.Args[2:])
	case "-h", "-help", "--help":
		usage()
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `natpunch: a minimal UDP NAT hole-punching demo

Usage:
  natpunch server  -listen :9000 -udp :9001
  natpunch punch   -server <host:9000> -udp-server <host:9001>
                   -session <room> -name <alice|bob> -local :0 [-text "hello"]

Run 'natpunch server -h' or 'natpunch punch -h' for flag details.`)
}

// =====================================================================
// Signaling server
// =====================================================================

type peer struct {
	name string
	conn net.Conn // TCP signaling conn
	addr *net.UDPAddr

	// out is a per-peer FIFO of outgoing signaling messages. A single
	// writer goroutine drains `out` and JSON-encodes each Msg onto conn.
	// This guarantees that any message enqueued before the peer was
	// visible to matchLocked (i.e., the initial "registered") reaches the
	// peer before any subsequent "peer" messages from other goroutines.
	out chan Msg
	// done is closed when the writer goroutine has finished.
	done chan struct{}
	// closeOnce guards close(out) so concurrent cleanups don't panic.
	closeOnce sync.Once
}

// newPeer wires up a peer's writer goroutine. The writer reads from out
// until the channel is closed, encoding each Msg as a JSON line.
func newPeer(name string, conn net.Conn) *peer {
	p := &peer{
		name: name,
		conn: conn,
		out:  make(chan Msg, 16),
		done: make(chan struct{}),
	}
	go p.writer()
	return p
}

func (p *peer) writer() {
	defer close(p.done)
	enc := json.NewEncoder(p.conn)
	for m := range p.out {
		if err := enc.Encode(m); err != nil {
			return
		}
	}
}

// send enqueues m for delivery. Returns false if the channel is closed
// or the buffer is full (callers should treat both as "the peer is gone").
func (p *peer) send(m Msg) bool {
	defer func() { _ = recover() }() // channel closed concurrently
	select {
	case p.out <- m:
		return true
	default:
		return false
	}
}

// shutdown stops the writer goroutine and waits for it to exit.
func (p *peer) shutdown() {
	p.closeOnce.Do(func() { close(p.out) })
	<-p.done
}

type session struct {
	mu      chan struct{} // tiny mutex stand-in; capacity 1
	peers   map[string]*peer
	pending map[string]*net.UDPAddr // probes seen before TCP register
}

type server struct {
	tcp, udp string

	mu       chan struct{} // capacity 1
	sessions map[string]*session
}

func runServer(args []string) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	tcpAddr := fs.String("listen", ":9000", "TCP listen address for signaling")
	udpAddr := fs.String("udp", ":9001", "UDP listen address for peer probes")
	_ = fs.Parse(args)

	s := &server{
		tcp: *tcpAddr, udp: *udpAddr,
		mu:       make(chan struct{}, 1),
		sessions: map[string]*session{},
	}
	if err := s.run(); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func (s *server) lock()   { s.mu <- struct{}{} }
func (s *server) unlock() { <-s.mu }

func (s *server) sess(name string) *session {
	s.lock()
	defer s.unlock()
	ss, ok := s.sessions[name]
	if !ok {
		ss = &session{
			mu:      make(chan struct{}, 1),
			peers:   map[string]*peer{},
			pending: map[string]*net.UDPAddr{},
		}
		s.sessions[name] = ss
	}
	return ss
}

func (s *server) run() error {
	uAddr, err := net.ResolveUDPAddr("udp", s.udp)
	if err != nil {
		return err
	}
	u, err := net.ListenUDP("udp", uAddr)
	if err != nil {
		return err
	}
	defer u.Close()
	log.Printf("server: UDP probe listener on %s", u.LocalAddr())
	go s.readProbes(u)

	ln, err := net.Listen("tcp", s.tcp)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Printf("server: TCP signaling on %s", ln.Addr())

	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handle(c)
	}
}

// readProbes consumes UDP datagrams from peers. The payload format is
// "<session>/<name>"; the source address of the datagram is the peer's
// public UDP endpoint as seen by the server.
func (s *server) readProbes(c *net.UDPConn) {
	buf := make([]byte, 256)
	for {
		n, src, err := c.ReadFromUDP(buf)
		if err != nil {
			return
		}
		sessName, peerName, ok := splitProbe(string(buf[:n]))
		if !ok {
			continue
		}
		ss := s.sess(sessName)
		ss.mu <- struct{}{}
		if p, ok := ss.peers[peerName]; ok {
			p.addr = src
			log.Printf("server: %s/%s public UDP = %s", sessName, peerName, src)
			s.matchLocked(ss)
		} else {
			ss.pending[peerName] = src
			log.Printf("server: probe buffered for %s/%s from %s", sessName, peerName, src)
		}
		<-ss.mu
	}
}

func (s *server) handle(c net.Conn) {
	defer c.Close()
	dec := json.NewDecoder(c)

	var m Msg
	if err := dec.Decode(&m); err != nil {
		log.Printf("server: decode: %v", err)
		return
	}

	if m.Type != msgRegister {
		p := newPeer("", c)
		p.send(Msg{Type: msgError, Text: "first message must be register"})
		p.shutdown()
		return
	}

	p := newPeer(m.Name, c)

	// CRITICAL: send "registered" BEFORE we add the peer to ss.peers.
	// Once the peer is visible, any other goroutine running matchLocked
	// (e.g., the probe reader) could enqueue a "peer" message for it.
	// By enqueueing "registered" first, the writer's FIFO order
	// guarantees the client receives registered before any peer info.
	p.send(Msg{Type: msgRegistered})

	ss := s.sess(m.Session)
	ss.mu <- struct{}{}
	if pa, ok := ss.pending[m.Name]; ok {
		p.addr = pa
		delete(ss.pending, m.Name)
	}
	ss.peers[m.Name] = p
	log.Printf("server: registered %s/%s via %s", m.Session, m.Name, c.RemoteAddr())
	s.matchLocked(ss)
	<-ss.mu

	defer func() {
		ss.mu <- struct{}{}
		delete(ss.peers, m.Name)
		<-ss.mu
		p.shutdown()
		log.Printf("server: %s/%s gone", m.Session, m.Name)
	}()

	// Hold the connection open; the only purpose is to detect disconnect.
	for {
		var discard json.RawMessage
		if err := dec.Decode(&discard); err != nil {
			return
		}
	}
}

// matchLocked notifies every known peer (with a known public UDP addr) of
// every other known peer. Caller must hold ss.mu.
func (s *server) matchLocked(ss *session) {
	var known []*peer
	for _, p := range ss.peers {
		if p.addr != nil {
			known = append(known, p)
		}
	}
	if len(known) < 2 {
		return
	}
	for i, a := range known {
		for j, b := range known {
			if i == j {
				continue
			}
			m := Msg{Type: msgPeer, Name: b.name, Addr: b.addr.String()}
			if !a.send(m) {
				log.Printf("server: send queue full for %s", a.name)
			}
		}
	}
	log.Printf("server: matched %d peers in session", len(known))
}

func splitProbe(s string) (sess, name string, ok bool) {
	i := strings.IndexByte(s, '/')
	if i <= 0 || i == len(s)-1 {
		return "", "", false
	}
	return s[:i], s[i+1:], true
}

// =====================================================================
// Punch client
// =====================================================================

type client struct {
	serverTCP string
	serverUDP string
	session   string
	name      string
	localUDP  string
	text      string // if non-empty, send once after establishment and exit
}

func runPunch(args []string) {
	fs := flag.NewFlagSet("punch", flag.ExitOnError)
	server := fs.String("server", "127.0.0.1:9000", "signaling server TCP address")
	udpSrv := fs.String("udp-server", "127.0.0.1:9001", "signaling server UDP address (for probes)")
	session := fs.String("session", "demo", "session/room name")
	name := fs.String("name", "", "this peer's name (required)")
	local := fs.String("local", ":0", "local UDP bind address (:0 = pick ephemeral)")
	text := fs.String("text", "", "if set, send this text once the punch lands, then exit")
	_ = fs.Parse(args)
	if *name == "" {
		fmt.Fprintln(os.Stderr, "-name is required")
		os.Exit(2)
	}
	c := &client{
		serverTCP: *server, serverUDP: *udpSrv,
		session: *session, name: *name,
		localUDP: *local, text: *text,
	}
	if err := c.run(); err != nil {
		log.Fatalf("client %s: %v", c.name, err)
	}
}

// dialWithRetry tries to dial `addr` a few times before giving up. This
// is purely a usability helper — in production the server is presumed up
// and a single attempt is fine.
func dialWithRetry(network, addr string, total time.Duration) (net.Conn, error) {
	deadline := time.Now().Add(total)
	var lastErr error
	backoff := 100 * time.Millisecond
	for {
		c, err := net.Dial(network, addr)
		if err == nil {
			return c, nil
		}
		lastErr = err
		if time.Now().After(deadline) {
			return nil, lastErr
		}
		time.Sleep(backoff)
		if backoff < 1*time.Second {
			backoff *= 2
		}
	}
}

func (c *client) run() error {
	// 1. Open the UDP socket FIRST. This is the socket whose NAT mapping
	//    will be used for the hole punch — it must exist before we send
	//    the first packet.
	lAddr, err := net.ResolveUDPAddr("udp", c.localUDP)
	if err != nil {
		return fmt.Errorf("resolve local: %w", err)
	}
	udp, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		return fmt.Errorf("listen UDP: %w", err)
	}
	defer udp.Close()
	log.Printf("client[%s]: UDP socket on %s", c.name, udp.LocalAddr())

	// 2. Dial the signaling server over TCP. Retry briefly so the test
	//    demo (where the server may still be starting up) is forgiving.
	tcp, err := dialWithRetry("tcp", c.serverTCP, 3*time.Second)
	if err != nil {
		return fmt.Errorf("dial TCP: %w", err)
	}
	defer tcp.Close()
	log.Printf("client[%s]: TCP connected to %s", c.name, c.serverTCP)

	// 3. Send a UDP probe so the server learns our public endpoint. This
	//    may race with our TCP register; the server buffers probes.
	srvUDP, err := net.ResolveUDPAddr("udp", c.serverUDP)
	if err != nil {
		return err
	}
	probe := c.session + "/" + c.name
	if _, err := udp.WriteToUDP([]byte(probe), srvUDP); err != nil {
		// Non-fatal: server may not be listening on UDP yet, or our NAT
		// hasn't mapped our socket yet (rare). The TCP register path
		// will retry via matchLocked if we beat the probe.
		log.Printf("client[%s]: probe send (non-fatal): %v", c.name, err)
	}

	// 4. Register over TCP and wait for ack.
	enc := json.NewEncoder(tcp)
	dec := json.NewDecoder(tcp)
	if err := enc.Encode(Msg{Type: msgRegister, Session: c.session, Name: c.name}); err != nil {
		return err
	}
	var ack Msg
	if err := dec.Decode(&ack); err != nil {
		return err
	}
	if ack.Type != msgRegistered {
		return fmt.Errorf("register not acked: %+v", ack)
	}

	// 5. Wait for the server to introduce us to our peer.
	var pm Msg
	if err := dec.Decode(&pm); err != nil {
		return err
	}
	if pm.Type != msgPeer {
		return fmt.Errorf("expected peer msg, got %+v", pm)
	}
	peerAddr, err := net.ResolveUDPAddr("udp", pm.Addr)
	if err != nil {
		return err
	}
	log.Printf("client[%s]: peer=%s at %s — punching...", c.name, pm.Name, peerAddr)

	// 6. Hole-punch: keep firing PUNCH heartbeats AND reading. The first
	//    inbound datagram from the peer confirms bidirectional mapping.
	done := make(chan struct{})
	punchedOnce := make(chan struct{})

	// Heartbeat sender. Exits when `done` closes or write fails.
	go func() {
		t := time.NewTicker(200 * time.Millisecond)
		defer t.Stop()
		i := 0
		for {
			select {
			case <-done:
				return
			case <-t.C:
				i++
				pkt := fmt.Sprintf("PUNCH|%s|%d", c.name, i)
				if _, err := udp.WriteToUDP([]byte(pkt), peerAddr); err != nil {
					return
				}
				if i == 1 || i%50 == 0 {
					log.Printf("client[%s]: PUNCH #%d -> %s", c.name, i, peerAddr)
				}
			}
		}
	}()

	// Reader. Exits on UDP read error or `done`.
	go func() {
		defer close(done)
		buf := make([]byte, 1500)
		for {
			n, src, err := udp.ReadFromUDP(buf)
			if err != nil {
				log.Printf("client[%s]: read: %v", c.name, err)
				return
			}
			// First packet from anyone (the peer) closes punchedOnce.
			select {
			case <-punchedOnce:
			default:
				close(punchedOnce)
				log.Printf("client[%s]: NAT PUNCH ESTABLISHED with %s", c.name, src)
			}
			data := string(buf[:n])
			switch {
			case strings.HasPrefix(data, "PUNCH|"):
				// heartbeat, no log spam
			case strings.HasPrefix(data, "CHAT|"):
				rest := strings.TrimPrefix(data, "CHAT|")
				pipe := strings.IndexByte(rest, '|')
				if pipe < 0 {
					continue
				}
				fmt.Printf("[%s] %s\n", rest[:pipe], rest[pipe+1:])
			default:
				log.Printf("client[%s]: recv %q from %s", c.name, data, src)
			}
		}
	}()

	// 7. Wait for the punch to be established.
	select {
	case <-punchedOnce:
	case <-done:
		return fmt.Errorf("reader exited before any peer packet arrived")
	case <-time.After(30 * time.Second):
		return fmt.Errorf("punch timeout: no packet from peer in 30s")
	}

	// 8. One-shot mode: send -text and exit.
	if c.text != "" {
		msg := "CHAT|" + c.name + "|" + c.text
		if _, err := udp.WriteToUDP([]byte(msg), peerAddr); err != nil {
			return err
		}
		log.Printf("client[%s]: sent one-shot text; exiting after brief wait", c.name)
		time.Sleep(500 * time.Millisecond)
		return nil
	}

	// 9. Interactive mode: stdin -> peer, peer -> stdout.
	log.Printf("client[%s]: punch OK. type lines and press Enter to chat; Ctrl-D/Ctrl-C to quit.", c.name)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			msg := "CHAT|" + c.name + "|" + sc.Text()
			if _, err := udp.WriteToUDP([]byte(msg), peerAddr); err != nil {
				return
			}
		}
		// EOF on stdin: close UDP to break the reader goroutine.
		_ = udp.Close()
	}()

	<-done
	return nil
}
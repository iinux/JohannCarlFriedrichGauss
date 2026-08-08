// Command sd is a minimal SSH server: it authenticates a single
// username/password pair and gives the client an interactive shell (PTY)
// or runs a one-off exec command, all configurable via CLI flags.
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
	"sd/lanproxy"
)

func main() {
	host := flag.String("host", "0.0.0.0", "listen address")
	port := flag.Int("port", 2222, "listen port")
	username := flag.String("user", "root", "required login username")
	password := flag.String("pass", "", "required login password")
	hostKeyPath := flag.String("hostkey", "", "path to PEM private key used as host key; generated in memory if empty")
	shellPath := flag.String("shell", "", "shell binary to launch on login; defaults to $SHELL or /bin/sh")
	lanProxyClientKey := flag.String("k", "", "lp client key")
	lanProxyServerHost := flag.String("s", "", "lp server host")
	lanProxyServerPort := flag.Int("p", 4900, "lp server port")
	lanProxySSL := flag.Bool("ssl", false, "enable lp ssl")
	lanProxyCertPath := flag.String("cer", "", "lp ssl cert path; skips verify if empty")
	flag.Parse()

	go lanproxy.LanProxyMain(*lanProxyClientKey, *lanProxyServerHost, *lanProxyServerPort, *lanProxySSL, *lanProxyCertPath)

	if *password == "" {
		log.Fatal("must supply -pass")
	}

	signer, err := loadOrGenerateHostKey(*hostKeyPath)
	if err != nil {
		log.Fatalf("host key: %v", err)
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == *username && string(pass) == *password {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}
	config.AddHostKey(signer)

	addr := net.JoinHostPort(*host, strconv.Itoa(*port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
	log.Printf("sd listening on %s (user=%s)", addr, *username)

	sh := *shellPath
	if sh == "" {
		sh = os.Getenv("SHELL")
	}
	if sh == "" {
		sh = "/bin/sh"
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleConn(nConn, config, sh)
	}
}

func loadOrGenerateHostKey(path string) (ssh.Signer, error) {
	if path != "" {
		keyBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return ssh.ParsePrivateKey(keyBytes)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

func handleConn(nConn net.Conn, config *ssh.ServerConfig, shellPath string) {
	defer nConn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Printf("handshake failed from %s: %v", nConn.RemoteAddr(), err)
		return
	}
	defer sshConn.Close()
	log.Printf("login: user=%s from %s", sshConn.User(), sshConn.RemoteAddr())

	forwards := newRemoteForwardSet()
	defer forwards.closeAll()
	go handleGlobalRequests(sshConn, reqs, forwards)

	for newChannel := range chans {
		switch newChannel.ChannelType() {
		case "session":
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Printf("channel accept: %v", err)
				continue
			}
			go handleSession(channel, requests, shellPath, func() {
				sshConn.Close()
			})
		case "direct-tcpip":
			go handleDirectTCPIP(newChannel)
		default:
			newChannel.Reject(ssh.UnknownChannelType, "unsupported channel type")
		}
	}
}

type directTCPIPPayload struct {
	DestAddr       string
	DestPort       uint32
	OriginatorAddr string
	OriginatorPort uint32
}

func handleDirectTCPIP(newChannel ssh.NewChannel) {
	var payload directTCPIPPayload
	if err := ssh.Unmarshal(newChannel.ExtraData(), &payload); err != nil {
		newChannel.Reject(ssh.ConnectionFailed, "invalid direct-tcpip payload")
		return
	}

	addr := net.JoinHostPort(payload.DestAddr, strconv.FormatUint(uint64(payload.DestPort), 10))
	target, err := net.Dial("tcp", addr)
	if err != nil {
		newChannel.Reject(ssh.ConnectionFailed, err.Error())
		log.Printf("direct-tcpip dial %s failed: %v", addr, err)
		return
	}

	channel, requests, err := newChannel.Accept()
	if err != nil {
		target.Close()
		log.Printf("direct-tcpip accept: %v", err)
		return
	}
	go ssh.DiscardRequests(requests)
	log.Printf("direct-tcpip: %s:%d -> %s", payload.OriginatorAddr, payload.OriginatorPort, addr)
	proxyTCP(channel, target)
}

type tcpipForwardPayload struct {
	Addr string
	Port uint32
}

type forwardedTCPIPPayload struct {
	ConnectedAddr  string
	ConnectedPort  uint32
	OriginatorAddr string
	OriginatorPort uint32
}

type remoteForwardSet struct {
	mu        sync.Mutex
	listeners map[string]net.Listener
}

func newRemoteForwardSet() *remoteForwardSet {
	return &remoteForwardSet{listeners: make(map[string]net.Listener)}
}

func (f *remoteForwardSet) add(addr string, listener net.Listener) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.listeners[addr]; exists {
		return false
	}
	f.listeners[addr] = listener
	return true
}

func (f *remoteForwardSet) remove(addr string) bool {
	f.mu.Lock()
	listener, exists := f.listeners[addr]
	if exists {
		delete(f.listeners, addr)
	}
	f.mu.Unlock()
	if exists {
		listener.Close()
	}
	return exists
}

func (f *remoteForwardSet) closeAll() {
	f.mu.Lock()
	listeners := f.listeners
	f.listeners = make(map[string]net.Listener)
	f.mu.Unlock()
	for _, listener := range listeners {
		listener.Close()
	}
}

func handleGlobalRequests(conn *ssh.ServerConn, reqs <-chan *ssh.Request, forwards *remoteForwardSet) {
	for req := range reqs {
		switch req.Type {
		case "tcpip-forward":
			handleTCPIPForward(conn, req, forwards)
		case "cancel-tcpip-forward":
			handleCancelTCPIPForward(req, forwards)
		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}
}

func handleTCPIPForward(conn *ssh.ServerConn, req *ssh.Request, forwards *remoteForwardSet) {
	var payload tcpipForwardPayload
	if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
		req.Reply(false, nil)
		return
	}

	addr := net.JoinHostPort(payload.Addr, strconv.FormatUint(uint64(payload.Port), 10))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("tcpip-forward listen %s failed: %v", addr, err)
		req.Reply(false, nil)
		return
	}

	actualAddr := listener.Addr().(*net.TCPAddr)
	actualPort := uint32(actualAddr.Port)
	forwardKey := net.JoinHostPort(payload.Addr, strconv.Itoa(actualAddr.Port))
	if !forwards.add(forwardKey, listener) {
		listener.Close()
		req.Reply(false, nil)
		return
	}

	if payload.Port == 0 {
		req.Reply(true, ssh.Marshal(struct{ Port uint32 }{Port: actualPort}))
	} else {
		req.Reply(true, nil)
	}
	log.Printf("tcpip-forward listening on %s", listener.Addr())
	go acceptRemoteForward(conn, listener, forwards, forwardKey, payload.Addr, actualPort)
}

func handleCancelTCPIPForward(req *ssh.Request, forwards *remoteForwardSet) {
	var payload tcpipForwardPayload
	if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
		req.Reply(false, nil)
		return
	}

	addr := net.JoinHostPort(payload.Addr, strconv.FormatUint(uint64(payload.Port), 10))
	req.Reply(forwards.remove(addr), nil)
}

func acceptRemoteForward(conn *ssh.ServerConn, listener net.Listener, forwards *remoteForwardSet, forwardKey string, connectedHost string, connectedPort uint32) {
	defer forwards.remove(forwardKey)
	for {
		inbound, err := listener.Accept()
		if err != nil {
			return
		}
		go openForwardedTCPIP(conn, inbound, connectedHost, connectedPort)
	}
}

func openForwardedTCPIP(conn *ssh.ServerConn, inbound net.Conn, connectedHost string, connectedPort uint32) {
	defer inbound.Close()

	originHost, originPort := splitTCPAddr(inbound.RemoteAddr())
	payload := forwardedTCPIPPayload{
		ConnectedAddr:  connectedHost,
		ConnectedPort:  connectedPort,
		OriginatorAddr: originHost,
		OriginatorPort: originPort,
	}
	channel, requests, err := conn.OpenChannel("forwarded-tcpip", ssh.Marshal(payload))
	if err != nil {
		log.Printf("forwarded-tcpip open failed: %v", err)
		return
	}
	go ssh.DiscardRequests(requests)
	proxyTCP(channel, inbound)
}

func splitTCPAddr(addr net.Addr) (string, uint32) {
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		return tcpAddr.IP.String(), uint32(tcpAddr.Port)
	}
	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return addr.String(), 0
	}
	n, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return host, 0
	}
	return host, uint32(n)
}

func proxyTCP(channel ssh.Channel, conn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(channel, conn)
		channel.CloseWrite()
	}()

	go func() {
		defer wg.Done()
		io.Copy(conn, channel)
		closeWrite(conn)
	}()

	wg.Wait()
	channel.Close()
	conn.Close()
}

func closeWrite(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
		return
	}
	conn.Close()
}

func handleSession(channel ssh.Channel, requests <-chan *ssh.Request, shellPath string, onDone func()) {
	defer channel.Close()
	defer onDone()

	var (
		ptmx     *os.File
		cmd      *exec.Cmd
		termOnce sync.Once
	)
	closePty := func() {
		termOnce.Do(func() {
			if ptmx != nil {
				ptmx.Close()
			}
		})
	}
	defer closePty()

	for req := range requests {
		switch req.Type {
		case "pty-req":
			cols, rows, ok := parsePtyReq(req.Payload)
			if !ok {
				req.Reply(false, nil)
				continue
			}
			cmd = exec.Command(shellPath)
			cmd.Env = append(os.Environ(), "TERM=xterm")
			f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
			if err != nil {
				log.Printf("pty start: %v", err)
				req.Reply(false, nil)
				continue
			}
			ptmx = f
			req.Reply(true, nil)

		case "window-change":
			cols, rows, ok := parseWinChange(req.Payload)
			if ok && ptmx != nil {
				pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
			}

		case "shell":
			req.Reply(true, nil)
			if ptmx == nil {
				// No pty was requested; run the shell without one.
				cmd = exec.Command(shellPath)
				cmd.Env = os.Environ()
				runWithoutPty(cmd, channel)
				return
			}
			go func() {
				io.Copy(ptmx, channel)
			}()
			io.Copy(channel, ptmx)
			cmd.Wait()
			return

		case "exec":
			command, ok := parseExecReq(req.Payload)
			if !ok {
				req.Reply(false, nil)
				continue
			}
			req.Reply(true, nil)
			c := exec.Command(shellPath, "-c", command)
			c.Env = os.Environ()
			runWithoutPty(c, channel)
			return

		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}
}

func runWithoutPty(cmd *exec.Cmd, channel ssh.Channel) {
	stdin, _ := cmd.StdinPipe()
	cmd.Stdout = channel
	cmd.Stderr = channel.Stderr()

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(channel.Stderr(), "failed to start: %v\n", err)
		exitStatus(channel, 1)
		return
	}
	go func() {
		io.Copy(stdin, channel)
		stdin.Close()
	}()

	err := cmd.Wait()
	code := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		code = 1
	}
	exitStatus(channel, uint32(code))
}

func exitStatus(channel ssh.Channel, code uint32) {
	var payload struct{ Status uint32 }
	payload.Status = code
	channel.SendRequest("exit-status", false, ssh.Marshal(&payload))
}

// parsePtyReq decodes the pty-req payload (RFC 4254 6.2): TERM string then
// terminal dimensions; we only need columns/rows.
func parsePtyReq(payload []byte) (cols, rows int, ok bool) {
	if len(payload) < 4 {
		return 0, 0, false
	}
	termLen := int(payload[0])<<24 | int(payload[1])<<16 | int(payload[2])<<8 | int(payload[3])
	i := 4 + termLen
	if len(payload) < i+8 {
		return 0, 0, false
	}
	cols = int(payload[i])<<24 | int(payload[i+1])<<16 | int(payload[i+2])<<8 | int(payload[i+3])
	rows = int(payload[i+4])<<24 | int(payload[i+5])<<16 | int(payload[i+6])<<8 | int(payload[i+7])
	return cols, rows, true
}

func parseWinChange(payload []byte) (cols, rows int, ok bool) {
	if len(payload) < 8 {
		return 0, 0, false
	}
	cols = int(payload[0])<<24 | int(payload[1])<<16 | int(payload[2])<<8 | int(payload[3])
	rows = int(payload[4])<<24 | int(payload[5])<<16 | int(payload[6])<<8 | int(payload[7])
	return cols, rows, true
}

func parseExecReq(payload []byte) (string, bool) {
	if len(payload) < 4 {
		return "", false
	}
	n := int(payload[0])<<24 | int(payload[1])<<16 | int(payload[2])<<8 | int(payload[3])
	if len(payload) < 4+n {
		return "", false
	}
	return string(payload[4 : 4+n]), true
}

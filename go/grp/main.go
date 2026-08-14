package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	handshakeTimeout = 10 * time.Second
	requestTimeout   = 15 * time.Second
	heartbeatEvery   = 10 * time.Second
	controlTimeout   = 35 * time.Second
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	if len(os.Args) < 2 {
		usage()
	}

	var err error
	switch os.Args[1] {
	case "-s":
		if len(os.Args) != 4 {
			usage()
		}
		err = runServer(os.Args[2], os.Args[3])
	case "-c":
		if len(os.Args) != 6 {
			usage()
		}
		err = runClient(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	default:
		usage()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
    //fmt.Fprintf(os.Stderr, "用法:\n  %s -s <控制端口> <请求端口>\n  %s -c <服务端地址> <服务端端口> <目标地址> <目标端口>\n", os.Args[0], os.Args[0])
	os.Exit(2)
}

func listenPort(port string) (net.Listener, error) {
	if _, err := validPort(port); err != nil {
		return nil, err
	}
	return net.Listen("tcp", ":"+port)
}

func validPort(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return 0, fmt.Errorf("无效端口 %q", port)
	}
	return p, nil
}

type controlConn struct {
	conn net.Conn
	mu   sync.Mutex
}

func (c *controlConn) send(line string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(handshakeTimeout))
	_, err := io.WriteString(c.conn, line+"\n")
	return err
}

type server struct {
	mu      sync.Mutex
	control *controlConn
	pending map[string]chan net.Conn
}

func runServer(controlPort, requestPort string) error {
	controlListener, err := listenPort(controlPort)
	if err != nil {
		return fmt.Errorf("监听控制端口: %w", err)
	}
	defer controlListener.Close()
	requestListener, err := listenPort(requestPort)
	if err != nil {
		return fmt.Errorf("监听请求端口: %w", err)
	}
	defer requestListener.Close()

	s := &server{pending: make(map[string]chan net.Conn)}
	//log.Printf("服务端已启动：控制端口 %s，请求端口 %s", controlPort, requestPort)
	go func() {
		for {
			conn, acceptErr := controlListener.Accept()
			if acceptErr != nil {
				return
			}
			go s.acceptInternal(conn)
		}
	}()

	for {
		conn, err := requestListener.Accept()
		if err != nil {
			return err
		}
		go s.handleUser(conn)
	}
}

func (s *server) acceptInternal(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(handshakeTimeout))
	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')
	if err != nil {
		conn.Close()
		return
	}
	fields := strings.Fields(line)
	if len(fields) == 1 && fields[0] == "CONTROL" {
		s.installControl(conn, r)
		return
	}
	if len(fields) == 2 && fields[0] == "DATA" {
		_ = conn.SetReadDeadline(time.Time{})
		s.deliverData(fields[1], &bufferedConn{Conn: conn, r: r})
		return
	}
	conn.Close()
}

func (s *server) installControl(conn net.Conn, r *bufio.Reader) {
	c := &controlConn{conn: conn}
	s.mu.Lock()
	old := s.control
	s.control = c
	s.mu.Unlock()
	if old != nil {
		old.conn.Close()
	}
	log.Printf("控制连接已建立：%s", conn.RemoteAddr())

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(heartbeatEvery)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if c.send("PING") != nil {
					conn.Close()
					return
				}
			case <-done:
				return
			}
		}
	}()
	for {
		_ = conn.SetReadDeadline(time.Now().Add(controlTimeout))
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) != "PONG" {
			break
		}
	}
	close(done)
	conn.Close()
	s.mu.Lock()
	if s.control == c {
		s.control = nil
	}
	s.mu.Unlock()
	log.Printf("控制连接已断开：%s", conn.RemoteAddr())
}

func (s *server) handleUser(user net.Conn) {
	id, err := randomID()
	if err != nil {
		user.Close()
		return
	}
	// An unbuffered channel ensures a late data connection is closed instead of
	// being left queued after the user-side request has timed out.
	ch := make(chan net.Conn)
	s.mu.Lock()
	s.pending[id] = ch
	c := s.control
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.pending, id)
		s.mu.Unlock()
	}()
	if c == nil || c.send("OPEN "+id) != nil {
		user.Close()
		return
	}

	select {
	case data := <-ch:
		bridge(user, data)
	case <-time.After(requestTimeout):
		log.Printf("等待数据连接超时：%s", id)
		user.Close()
	}
}

func (s *server) deliverData(id string, conn net.Conn) {
	s.mu.Lock()
	ch := s.pending[id]
	s.mu.Unlock()
	if ch == nil {
		conn.Close()
		return
	}
	select {
	case ch <- conn:
	default:
		conn.Close()
	}
}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func runClient(serverHost, serverPort, targetHost, targetPort string) error {
	if _, err := validPort(serverPort); err != nil {
		return err
	}
	if _, err := validPort(targetPort); err != nil {
		return err
	}
	serverAddr := net.JoinHostPort(serverHost, serverPort)
	targetAddr := net.JoinHostPort(targetHost, targetPort)
    //log.Printf("客户端已启动：服务端 %s，目标 %s", serverAddr, targetAddr)
	for {
		if err := clientSession(serverAddr, targetAddr); err != nil {
			log.Printf("控制连接断开：%v；3 秒后重连", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func clientSession(serverAddr, targetAddr string) error {
	conn, err := net.DialTimeout("tcp", serverAddr, handshakeTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := io.WriteString(conn, "CONTROL\n"); err != nil {
		return err
	}
    //log.Printf("控制连接已建立：%s", serverAddr)
	r := bufio.NewReader(conn)
	var writeMu sync.Mutex
	for {
		_ = conn.SetReadDeadline(time.Now().Add(controlTimeout))
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		fields := strings.Fields(line)
		if len(fields) == 1 && fields[0] == "PING" {
			writeMu.Lock()
			_ = conn.SetWriteDeadline(time.Now().Add(handshakeTimeout))
			_, err = io.WriteString(conn, "PONG\n")
			writeMu.Unlock()
			if err != nil {
				return err
			}
		} else if len(fields) == 2 && fields[0] == "OPEN" {
			go openTunnel(serverAddr, targetAddr, fields[1])
		} else {
			return errors.New("收到无效的控制消息")
		}
	}
}

func openTunnel(serverAddr, targetAddr, id string) {
	target, err := net.DialTimeout("tcp", targetAddr, requestTimeout)
	if err != nil {
		log.Printf("连接目标 %s 失败：%v", targetAddr, err)
		return
	}
	data, err := net.DialTimeout("tcp", serverAddr, requestTimeout)
	if err != nil {
		target.Close()
		log.Printf("建立数据连接失败：%v", err)
		return
	}
	_ = data.SetWriteDeadline(time.Now().Add(handshakeTimeout))
	if _, err = io.WriteString(data, "DATA "+id+"\n"); err != nil {
		target.Close()
		data.Close()
		return
	}
	_ = data.SetWriteDeadline(time.Time{})
	bridge(target, data)
}

type bufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *bufferedConn) Read(p []byte) (int, error) { return c.r.Read(p) }

func (c *bufferedConn) CloseWrite() error {
	if tcp, ok := c.Conn.(*net.TCPConn); ok {
		return tcp.CloseWrite()
	}
	return nil
}

func bridge(a, b net.Conn) {
	defer a.Close()
	defer b.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	copyOneWay := func(dst, src net.Conn) {
		defer wg.Done()
		_, _ = io.Copy(dst, src)
		if halfCloser, ok := dst.(interface{ CloseWrite() error }); ok {
			_ = halfCloser.CloseWrite()
		}
	}
	go copyOneWay(a, b)
	go copyOneWay(b, a)
	wg.Wait()
}

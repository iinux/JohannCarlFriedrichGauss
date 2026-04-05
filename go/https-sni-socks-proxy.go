package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ListenAddr   = ":443"
	DNSAddr      = ":53"           // DNS监听地址
	SOCKSAddr    = "127.0.0.1:1080" // SOCKS5代理地址
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
)

// DNSCache DNS缓存
type DNSCache struct {
	mu    sync.RWMutex
	cache map[string]string // domain -> ip
}

// NewDNSCache 创建新的DNS缓存
func NewDNSCache() *DNSCache {
	return &DNSCache{
		cache: make(map[string]string),
	}
}

// Get 获取缓存的IP
func (c *DNSCache) Get(domain string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ip, ok := c.cache[domain]
	return ip, ok
}

// Set 设置缓存的IP
func (c *DNSCache) Set(domain, ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[domain] = ip
}

// HTTPSProxyServer HTTPS代理服务器
type HTTPSProxyServer struct {
	listenAddr string
	socksAddr  string
	dnsCache   *DNSCache
}

// NewHTTPSProxyServer 创建新的HTTPS代理服务器
func NewHTTPSProxyServer(listenAddr, socksAddr string, dnsCache *DNSCache) *HTTPSProxyServer {
	return &HTTPSProxyServer{
		listenAddr: listenAddr,
		socksAddr:  socksAddr,
		dnsCache:   dnsCache,
	}
}

// Start 启动HTTPS代理服务器
func (s *HTTPSProxyServer) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("监听失败: %v", err)
	}
	defer listener.Close()

	fmt.Printf("HTTPS代理服务器已启动，监听地址: %s\n", s.listenAddr)
	fmt.Printf("SOCKS5代理地址: %s\n", s.socksAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection 处理客户端连接
func (s *HTTPSProxyServer) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// 设置读取超时
	clientConn.SetReadDeadline(time.Now().Add(ReadTimeout))

	// 读取TLS Client Hello消息，获取SNI
	serverName, clientHelloData, err := s.parseClientHello(clientConn)
	if err != nil {
		fmt.Printf("解析Client Hello失败: %v\n", err)
		return
	}

	if serverName == "" {
		fmt.Println("无法获取Server Name (SNI)")
		return
	}

	fmt.Printf("收到HTTPS连接，Server Name (SNI): %s\n", serverName)

	// 目标地址，默认443端口
	targetHost := serverName + ":443"

	// 通过SOCKS5代理连接到目标服务器
	targetConn, err := s.connectViaSOCKS(targetHost)
	if err != nil {
		fmt.Printf("SOCKS5连接失败: %v\n", err)
		return
	}
	defer targetConn.Close()

	fmt.Printf("已连接到目标服务器: %s，开始转发数据...\n", targetHost)

	// 首先将已读取的Client Hello数据转发到目标服务器
	if _, err := targetConn.Write(clientHelloData); err != nil {
		fmt.Printf("转发Client Hello失败: %v\n", err)
		return
	}

	// 双向转发数据
	s.relay(clientConn, targetConn)
}

// parseClientHello 从TLS Client Hello消息中解析SNI
func (s *HTTPSProxyServer) parseClientHello(conn net.Conn) (string, []byte, error) {
	// 读取TLS记录层头部 (5字节)
	header := make([]byte, 5)
	if _, err := io.ReadFull(conn, header); err != nil {
		return "", nil, fmt.Errorf("读取TLS头部失败: %v", err)
	}

	// 检查Content Type (应该是22 = Handshake)
	if header[0] != 0x16 {
		return "", nil, fmt.Errorf("不是TLS握手消息: %d", header[0])
	}

	// 获取记录层长度 (大端序)
	recordLen := int(header[3])<<8 | int(header[4])

	// 读取完整的握手消息
	handshake := make([]byte, recordLen)
	if _, err := io.ReadFull(conn, handshake); err != nil {
		return "", nil, fmt.Errorf("读取握手消息失败: %v", err)
	}

	// 检查Handshake Type (应该是1 = Client Hello)
	if len(handshake) < 1 || handshake[0] != 0x01 {
		return "", nil, fmt.Errorf("不是Client Hello消息: %d", handshake[0])
	}

	// 解析Client Hello
	serverName, err := s.extractSNI(handshake)
	if err != nil {
		return "", nil, err
	}

	// 组合完整的Client Hello数据（头部+握手消息）
	fullClientHello := append(header, handshake...)

	return serverName, fullClientHello, nil
}

// extractSNI 从Client Hello消息中提取SNI
func (s *HTTPSProxyServer) extractSNI(handshake []byte) (string, error) {
	if len(handshake) < 4 {
		return "", fmt.Errorf("Client Hello消息太短")
	}

	// 跳过Handshake Type (1字节) 和 Length (3字节)
	pos := 4

	// 跳过Client Version (2字节)
	if pos+2 > len(handshake) {
		return "", fmt.Errorf("Client Hello格式错误")
	}
	pos += 2

	// 跳过Random (32字节)
	if pos+32 > len(handshake) {
		return "", fmt.Errorf("Client Hello格式错误")
	}
	pos += 32

	// 跳过Session ID
	if pos+1 > len(handshake) {
		return "", fmt.Errorf("Client Hello格式错误")
	}
	sessionIDLen := int(handshake[pos])
	pos += 1 + sessionIDLen

	// 跳过Cipher Suites
	if pos+2 > len(handshake) {
		return "", fmt.Errorf("Client Hello格式错误")
	}
	cipherSuitesLen := int(handshake[pos])<<8 | int(handshake[pos+1])
	pos += 2 + cipherSuitesLen

	// 跳过Compression Methods
	if pos+1 > len(handshake) {
		return "", fmt.Errorf("Client Hello格式错误")
	}
	compressionMethodsLen := int(handshake[pos])
	pos += 1 + compressionMethodsLen

	// 解析Extensions
	if pos+2 > len(handshake) {
		return "", fmt.Errorf("没有Extensions")
	}
	extensionsLen := int(handshake[pos])<<8 | int(handshake[pos+1])
	pos += 2

	end := pos + extensionsLen
	if end > len(handshake) {
		return "", fmt.Errorf("Extensions长度错误")
	}

	// 遍历Extensions查找SNI (type = 0x0000)
	for pos < end {
		if pos+4 > end {
			break
		}

		extType := int(handshake[pos])<<8 | int(handshake[pos+1])
		extLen := int(handshake[pos+2])<<8 | int(handshake[pos+3])
		pos += 4

		if pos+extLen > end {
			break
		}

		if extType == 0x0000 { // SNI Extension
			return s.parseSNIExtension(handshake[pos : pos+extLen])
		}

		pos += extLen
	}

	return "", nil // 没有找到SNI
}

// parseSNIExtension 解析SNI Extension
func (s *HTTPSProxyServer) parseSNIExtension(data []byte) (string, error) {
	if len(data) < 2 {
		return "", fmt.Errorf("SNI Extension太短")
	}

	// SNI List Length
	sniListLen := int(data[0])<<8 | int(data[1])
	if sniListLen+2 > len(data) {
		return "", fmt.Errorf("SNI List长度错误")
	}

	pos := 2
	end := 2 + sniListLen

	for pos < end {
		if pos+3 > end {
			break
		}

		nameType := data[pos]
		nameLen := int(data[pos+1])<<8 | int(data[pos+2])
		pos += 3

		if pos+nameLen > end {
			break
		}

		if nameType == 0x00 { // host_name
			return string(data[pos : pos+nameLen]), nil
		}

		pos += nameLen
	}

	return "", fmt.Errorf("SNI中没有host_name")
}

// connectViaSOCKS 通过SOCKS5代理连接到目标服务器
func (s *HTTPSProxyServer) connectViaSOCKS(targetHost string) (net.Conn, error) {
	// 连接到SOCKS5代理
	socksConn, err := net.Dial("tcp", s.socksAddr)
	if err != nil {
		return nil, fmt.Errorf("连接SOCKS5代理失败: %v", err)
	}

	// SOCKS5握手
	if err := s.socks5Handshake(socksConn); err != nil {
		socksConn.Close()
		return nil, err
	}

	// 请求连接到目标服务器
	if err := s.socks5Connect(socksConn, targetHost); err != nil {
		socksConn.Close()
		return nil, err
	}

	return socksConn, nil
}

// socks5Handshake 执行SOCKS5握手
func (s *HTTPSProxyServer) socks5Handshake(conn net.Conn) error {
	// 发送握手请求：版本(5) + 认证方法数量(1) + 无认证(0)
	handshake := []byte{0x05, 0x01, 0x00}
	if _, err := conn.Write(handshake); err != nil {
		return fmt.Errorf("SOCKS5握手失败: %v", err)
	}

	// 读取响应
	response := make([]byte, 2)
	if _, err := io.ReadFull(conn, response); err != nil {
		return fmt.Errorf("SOCKS5握手响应失败: %v", err)
	}

	if response[0] != 0x05 {
		return fmt.Errorf("不支持的SOCKS版本: %d", response[0])
	}

	if response[1] != 0x00 {
		return fmt.Errorf("SOCKS5认证失败: %d", response[1])
	}

	return nil
}

// socks5Connect 请求SOCKS5代理连接到目标服务器
func (s *HTTPSProxyServer) socks5Connect(conn net.Conn, targetHost string) error {
	// 解析目标地址
	host, portStr, err := net.SplitHostPort(targetHost)
	if err != nil {
		return fmt.Errorf("解析目标地址失败: %v", err)
	}

	port := 0
	fmt.Sscanf(portStr, "%d", &port)

	// 构建CONNECT请求
	request := []byte{0x05, 0x01, 0x00} // VER, CMD(CONNECT), RSV

	// 判断是IPv4、IPv6还是域名
	ip := net.ParseIP(host)
	if ip4 := ip.To4(); ip4 != nil {
		// IPv4
		request = append(request, 0x01) // ATYP(IPv4)
		request = append(request, ip4...)
	} else if ip6 := ip.To16(); ip6 != nil {
		// IPv6
		request = append(request, 0x04) // ATYP(IPv6)
		request = append(request, ip6...)
	} else {
		// 域名
		request = append(request, 0x03) // ATYP(DOMAIN)
		request = append(request, byte(len(host)))
		request = append(request, []byte(host)...)
	}

	// 添加端口（大端序）
	request = append(request, byte(port>>8), byte(port&0xff))

	// 发送CONNECT请求
	if _, err := conn.Write(request); err != nil {
		return fmt.Errorf("SOCKS5 CONNECT请求失败: %v", err)
	}

	// 读取响应
	response := make([]byte, 4)
	if _, err := io.ReadFull(conn, response); err != nil {
		return fmt.Errorf("SOCKS5 CONNECT响应失败: %v", err)
	}

	if response[0] != 0x05 {
		return fmt.Errorf("不支持的SOCKS版本: %d", response[0])
	}

	if response[1] != 0x00 {
		return fmt.Errorf("SOCKS5 CONNECT失败，错误码: %d", response[1])
	}

	// 读取绑定的地址（跳过）
	switch response[3] {
	case 0x01: // IPv4
		addr := make([]byte, 4+2)
		io.ReadFull(conn, addr)
	case 0x03: // 域名
		lenBuf := make([]byte, 1)
		io.ReadFull(conn, lenBuf)
		addr := make([]byte, int(lenBuf[0])+2)
		io.ReadFull(conn, addr)
	case 0x04: // IPv6
		addr := make([]byte, 16+2)
		io.ReadFull(conn, addr)
	}

	return nil
}

// relay 双向转发数据
func (s *HTTPSProxyServer) relay(conn1, conn2 net.Conn) {
	var wg sync.WaitGroup
    wg.Add(2)

	go func() {
        defer wg.Done()
		io.Copy(conn1, conn2)
	}()

	go func() {
        defer wg.Done()
		io.Copy(conn2, conn1)
	}()

    wg.Wait()
}

// DNSServer DNS服务器
type DNSServer struct {
	addr     string
	socksAddr string
	cache    *DNSCache
	upstream string
}

// NewDNSServer 创建新的DNS服务器
func NewDNSServer(addr, socksAddr string, cache *DNSCache) *DNSServer {
	return &DNSServer{
		addr:      addr,
		socksAddr: socksAddr,
		cache:     cache,
		//upstream:  "8.8.8.8:53", // Google DNS
		upstream:  "223.5.5.5:53", // Google DNS
	}
}

// Start 启动DNS服务器
func (s *DNSServer) Start() error {
	// 监听UDP
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return fmt.Errorf("解析DNS地址失败: %v", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("监听DNS UDP失败: %v", err)
	}
	defer udpConn.Close()

	// 监听TCP
	tcpListener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("监听DNS TCP失败: %v", err)
	}
	defer tcpListener.Close()

	fmt.Printf("DNS服务器已启动，监听地址: %s\n", s.addr)

	// 启动UDP处理
	go s.handleUDP(udpConn)

	// 处理TCP请求
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Printf("接受DNS TCP连接失败: %v\n", err)
			continue
		}
		go s.handleTCP(conn)
	}
}

// handleUDP 处理DNS UDP请求
func (s *DNSServer) handleUDP(conn *net.UDPConn) {
	buffer := make([]byte, 4096)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("读取DNS UDP请求失败: %v\n", err)
			continue
		}

		go func(data []byte, addr *net.UDPAddr) {
			// 解析域名
			domain := s.parseDomain(data)
			if domain != "" {
				fmt.Printf("DNS查询(UDP): %s\n", domain)
			}

			// 检查是否为特定域名，返回127.0.0.1
			if s.isSpecialDomain(domain) {
				fmt.Printf("DNS特殊域名(UDP): %s -> 127.0.0.1\n", domain)
				response := s.buildResponse(data, "127.0.0.1")
				conn.WriteToUDP(response, addr)
				return
			}

			// 检查缓存
			if ip, ok := s.cache.Get(domain); ok {
				fmt.Printf("DNS缓存命中: %s -> %s\n", domain, ip)
				response := s.buildResponse(data, ip)
				conn.WriteToUDP(response, addr)
				return
			}

			// 转发到上游DNS
			response, err := s.forwardDNS(data)
			if err != nil {
				fmt.Printf("转发DNS失败: %v\n", err)
				return
			}

			// 解析响应并缓存
			if ip := s.parseIPFromResponse(response); ip != "" {
				s.cache.Set(domain, ip)
				fmt.Printf("DNS缓存: %s -> %s\n", domain, ip)
			}

			conn.WriteToUDP(response, addr)
		}(buffer[:n], clientAddr)
	}
}

// handleTCP 处理DNS TCP请求
func (s *DNSServer) handleTCP(conn net.Conn) {
	defer conn.Close()

	// 读取长度前缀（2字节，大端序）
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return
	}
	length := binary.BigEndian.Uint16(lenBuf)

	// 读取DNS请求
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return
	}

	// 解析域名
	domain := s.parseDomain(data)
	if domain != "" {
		fmt.Printf("DNS查询(TCP): %s\n", domain)
	}

	// 检查是否为特定域名，返回127.0.0.1
	if s.isSpecialDomain(domain) {
		fmt.Printf("DNS特殊域名(TCP): %s -> 127.0.0.1\n", domain)
		response := s.buildResponse(data, "127.0.0.1")
		s.sendTCPResponse(conn, response)
		return
	}

	// 检查缓存
	if ip, ok := s.cache.Get(domain); ok {
		fmt.Printf("DNS缓存命中: %s -> %s\n", domain, ip)
		response := s.buildResponse(data, ip)
		s.sendTCPResponse(conn, response)
		return
	}

	// 转发到上游DNS
	response, err := s.forwardDNS(data)
	if err != nil {
		fmt.Printf("转发DNS失败: %v\n", err)
		return
	}

	// 解析响应并缓存
	if ip := s.parseIPFromResponse(response); ip != "" {
		s.cache.Set(domain, ip)
		fmt.Printf("DNS缓存: %s -> %s\n", domain, ip)
	}

	s.sendTCPResponse(conn, response)
}

// sendTCPResponse 发送DNS TCP响应
func (s *DNSServer) sendTCPResponse(conn net.Conn, data []byte) {
	// TCP DNS响应需要长度前缀
	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, uint16(len(data)))
	conn.Write(lenBuf)
	conn.Write(data)
}

// parseDomain 从DNS请求中解析域名
func (s *DNSServer) parseDomain(data []byte) string {
	if len(data) < 12 {
		return ""
	}

	// 跳过DNS头部（12字节）
	pos := 12

	// 解析查询域名
	var domainParts []string
	for pos < len(data) {
		length := int(data[pos])
		if length == 0 {
			break
		}
		if pos+1+length > len(data) {
			return ""
		}
		domainParts = append(domainParts, string(data[pos+1:pos+1+length]))
		pos += 1 + length
	}

	return strings.Join(domainParts, ".")
}

// parseIPFromResponse 从DNS响应中解析IP地址
func (s *DNSServer) parseIPFromResponse(data []byte) string {
	if len(data) < 12 {
		return ""
	}

	// 获取回答数量
	ancount := binary.BigEndian.Uint16(data[6:8])
	if ancount == 0 {
		return ""
	}

	// 跳过DNS头部和查询部分
	pos := 12

	// 跳过查询域名
	for pos < len(data) {
		if data[pos] == 0 {
			pos++
			break
		}
		if data[pos]&0xC0 == 0xC0 {
			// 压缩指针
			pos += 2
			break
		}
		pos += 1 + int(data[pos])
	}

	// 跳过查询类型和类（4字节）
	pos += 4

	// 解析第一个回答
	if pos >= len(data) {
		return ""
	}

	// 跳过名称（可能是压缩指针）
	if data[pos]&0xC0 == 0xC0 {
		pos += 2
	} else {
		for pos < len(data) && data[pos] != 0 {
			pos += 1 + int(data[pos])
		}
		pos++
	}

	if pos+10 > len(data) {
		return ""
	}

	// 获取类型
	qtype := binary.BigEndian.Uint16(data[pos : pos+2])
	pos += 8 // 跳过类型、类、TTL

	// 获取数据长度
	rdlength := binary.BigEndian.Uint16(data[pos : pos+2])
	pos += 2

	if qtype == 1 && rdlength == 4 && pos+4 <= len(data) { // A记录
		ip := net.IP(data[pos : pos+4])
		return ip.String()
	}

	return ""
}

// buildResponse 构建DNS响应（从缓存）
func (s *DNSServer) buildResponse(request []byte, ip string) []byte {
	// 找到查询部分的末尾（包括域名、类型、类）
	pos := 12
	for pos < len(request) {
		if request[pos] == 0 {
			pos++
			break
		}
		// 处理压缩指针
		if request[pos]&0xC0 == 0xC0 {
			pos += 2
			break
		}
		pos += 1 + int(request[pos])
	}
	// 跳过QTYPE和QCLASS（4字节）
	queryEnd := pos + 4

	// 构建响应头部（12字节）
	response := make([]byte, 12)
	
	// 复制ID
	copy(response[0:2], request[0:2])
	
	// 设置响应标志
	response[2] = 0x81 // QR=1, OPCODE=0, AA=0, TC=0, RD=1
	response[3] = 0x80 // RA=1, Z=0, RCODE=0 (No Error)
	
	// 设置问题数量为1（大端序）
	response[4] = 0x00
	response[5] = 0x01
	
	// 设置回答数量为1（大端序）
	response[6] = 0x00
	response[7] = 0x01
	
	// NSCOUNT = 0
	response[8] = 0x00
	response[9] = 0x00
	
	// ARCOUNT = 0
	response[10] = 0x00
	response[11] = 0x00

	// 复制查询部分（原始问题）
	response = append(response, request[12:queryEnd]...)

	// 构建回答部分
	answer := []byte{
		0xC0, 0x0C, // 压缩指针指向查询域名（偏移12）
		0x00, 0x01, // 类型A
		0x00, 0x01, // 类IN
		0x00, 0x00, 0x0E, 0x10, // TTL = 3600秒
		0x00, 0x04, // 数据长度 = 4
	}

	// 添加IP地址
	ipBytes := net.ParseIP(ip).To4()
	answer = append(answer, ipBytes...)

	return append(response, answer...)
}

// forwardDNS 转发DNS请求到上游服务器
func (s *DNSServer) forwardDNS(data []byte) ([]byte, error) {
	// 连接到上游DNS服务器
	conn, err := net.Dial("udp", s.upstream)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 设置超时
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// 发送请求
	if _, err := conn.Write(data); err != nil {
		return nil, err
	}

	// 读取响应
	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}

// isSpecialDomain 检查是否为特殊域名
func (s *DNSServer) isSpecialDomain(domain string) bool {
	domain = strings.ToLower(domain)
	return strings.HasSuffix(domain, "cloudflarestorage.com") || 
	       strings.HasSuffix(domain, "docker.io")
}

func main() {
	// 解析命令行参数
	socksAddr := SOCKSAddr
	if len(os.Args) > 1 {
		socksAddr = os.Args[1]
	}

	listenAddr := ListenAddr
	if len(os.Args) > 2 {
		listenAddr = os.Args[2]
	}

	dnsAddr := DNSAddr
	if len(os.Args) > 3 {
		dnsAddr = os.Args[3]
	}

	// 创建共享的DNS缓存
	dnsCache := NewDNSCache()

	// 创建并启动DNS服务器
	dnsServer := NewDNSServer(dnsAddr, socksAddr, dnsCache)
	go func() {
		if err := dnsServer.Start(); err != nil {
			fmt.Printf("DNS服务器启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 创建并启动HTTPS代理服务器
	httpsServer := NewHTTPSProxyServer(listenAddr, socksAddr, dnsCache)

	fmt.Printf("启动HTTPS SNI代理服务器...\n")
	fmt.Printf("监听地址: %s\n", listenAddr)
	fmt.Printf("SOCKS5代理: %s\n", socksAddr)
	fmt.Printf("DNS监听: %s\n", dnsAddr)
	fmt.Println("使用方法: go run https-sni-socks-proxy.go [SOCKS5地址] [HTTPS监听地址] [DNS监听地址]")
	fmt.Println("示例: go run https-sni-socks-proxy.go 127.0.0.1:1080 :443 :53")

	if err := httpsServer.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}

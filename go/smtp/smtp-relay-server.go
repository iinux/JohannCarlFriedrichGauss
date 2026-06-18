package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	ListenAddr     = ":25"                    // SMTP监听地址
	TLSCertFile    = "server.crt"             // TLS证书文件
	TLSKeyFile     = "server.key"             // TLS私钥文件
	MaxMessageSize = 10 * 1024 * 1024         // 最大邮件大小 10MB
	ReadTimeout    = 60 * time.Second
	WriteTimeout   = 60 * time.Second
)

// SMTPConfig SMTP服务器配置
type SMTPConfig struct {
	ListenAddr     string
	TLSCertFile    string
	TLSKeyFile     string
	MaxMessageSize int64
	// 转发配置
	RelayHost     string // 转发SMTP服务器地址，如 smtp.gmail.com:587
	RelayUser     string // 转发邮箱用户名
	RelayPassword string // 转发邮箱密码
	RelayTo       string // 目标转发邮箱地址
}

// SMTPSession SMTP会话
type SMTPSession struct {
	conn       net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	config     *SMTPConfig
	helo       string
	from       string
	to         []string
	data       strings.Builder
	inData     bool
	tlsEnabled bool
}

// NewSMTPSession 创建新的SMTP会话
func NewSMTPSession(conn net.Conn, config *SMTPConfig) *SMTPSession {
	return &SMTPSession{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		config: config,
		to:     make([]string, 0),
	}
}

// WriteResponse 发送SMTP响应
func (s *SMTPSession) WriteResponse(code int, message string) error {
	_, err := fmt.Fprintf(s.writer, "%d %s\r\n", code, message)
	s.writer.Flush()
	return err
}

// ReadLine 读取一行
// 只去掉行尾 CRLF，保留前导空白
// Why: DATA 阶段的折叠头（RFC 5322）续行必须以空白开头，TrimSpace 会破坏头部结构
func (s *SMTPSession) ReadLine() (string, error) {
	line, err := s.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// Handle 处理SMTP会话
func (s *SMTPSession) Handle() {
	fmt.Printf("新连接: %s\n", s.conn.RemoteAddr().String())
	defer s.conn.Close()

	// 发送欢迎消息
	s.WriteResponse(220, "Go SMTP Relay Server Ready")

	for {
		s.conn.SetReadDeadline(time.Now().Add(ReadTimeout))
		line, err := s.ReadLine()
		if err != nil {
			fmt.Printf("读取命令失败: %v\n", err)
			return
		}

		if s.inData {
			s.handleData(line)
			continue
		}

		cmd, arg := s.parseCommand(line)
		fmt.Printf("收到命令: %s %s\n", cmd, arg)

		switch strings.ToUpper(cmd) {
		case "HELO", "EHLO":
			s.handleHelo(arg)
		case "MAIL":
			s.handleMail(arg)
		case "RCPT":
			s.handleRcpt(arg)
		case "DATA":
			s.handleDataStart()
		case "RSET":
			s.handleRset()
		case "VRFY":
			s.WriteResponse(252, "Cannot VRFY user")
		case "EXPN":
			s.WriteResponse(502, "Command not implemented")
		case "HELP":
			s.WriteResponse(214, "See RFC 5321")
		case "QUIT":
			s.WriteResponse(221, "Bye")
			return
		case "STARTTLS":
			s.handleStartTLS()
		case "AUTH":
			s.WriteResponse(502, "Command not implemented")
		default:
			s.WriteResponse(500, "Command not recognized")
		}
	}
}

// parseCommand 解析SMTP命令
func (s *SMTPSession) parseCommand(line string) (cmd, arg string) {
	parts := strings.SplitN(line, " ", 2)
	cmd = parts[0]
	if len(parts) > 1 {
		arg = parts[1]
	}
	return
}

// handleHelo 处理HELO/EHLO命令
func (s *SMTPSession) handleHelo(arg string) {
	s.helo = arg
	s.WriteResponse(250, "Hello "+arg)
}

// handleMail 处理MAIL FROM命令
func (s *SMTPSession) handleMail(arg string) {
	// 解析 MAIL FROM:<address>
	if !strings.HasPrefix(strings.ToUpper(arg), "FROM:") {
		s.WriteResponse(501, "Syntax error in parameters")
		return
	}

	addr := s.extractAddress(arg[5:])
	if addr == "" {
		s.WriteResponse(501, "Syntax error in parameters")
		return
	}

	s.from = addr
	s.WriteResponse(250, "OK")
}

// handleRcpt 处理RCPT TO命令
func (s *SMTPSession) handleRcpt(arg string) {
	// 解析 RCPT TO:<address>
	if !strings.HasPrefix(strings.ToUpper(arg), "TO:") {
		s.WriteResponse(501, "Syntax error in parameters")
		return
	}

	addr := s.extractAddress(arg[3:])
	if addr == "" {
		s.WriteResponse(501, "Syntax error in parameters")
		return
	}

	s.to = append(s.to, addr)
	s.WriteResponse(250, "OK")
}

// extractAddress 从 <address> 或 address 中提取邮箱地址
func (session *SMTPSession) extractAddress(str string) string {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "<") && strings.HasSuffix(str, ">") {
		return str[1 : len(str)-1]
	}
	return str
}

// handleDataStart 处理DATA命令开始
func (s *SMTPSession) handleDataStart() {
	if s.from == "" || len(s.to) == 0 {
		s.WriteResponse(503, "Bad sequence of commands")
		return
	}
	s.inData = true
	s.WriteResponse(354, "Start mail input; end with <CRLF>.<CRLF>")
}

// handleData 处理邮件数据
func (s *SMTPSession) handleData(line string) {
	if line == "." {
		// 邮件结束
		s.inData = false
		s.processMessage()
		s.resetTransaction()
		return
	}

	// 处理转义的点号 (.. -> .)
	if strings.HasPrefix(line, "..") {
		line = line[1:]
	}

	s.data.WriteString(line)
	s.data.WriteString("\r\n")

	// 检查邮件大小
	if int64(s.data.Len()) > s.config.MaxMessageSize {
		s.WriteResponse(552, "Message exceeds fixed maximum message size")
		s.inData = false
		s.resetTransaction()
	}
}

// handleRset 处理RSET命令
func (s *SMTPSession) handleRset() {
	s.resetTransaction()
	s.WriteResponse(250, "OK")
}

// resetTransaction 重置事务状态
func (s *SMTPSession) resetTransaction() {
	s.from = ""
	s.to = nil
	s.data.Reset()
	s.inData = false
}

// handleStartTLS 处理STARTTLS命令
func (s *SMTPSession) handleStartTLS() {
	if s.tlsEnabled {
		s.WriteResponse(503, "TLS already active")
		return
	}

	cert, err := tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
	if err != nil {
		fmt.Printf("加载TLS证书失败: %v\n", err)
		s.WriteResponse(454, "TLS not available")
		return
	}

	s.WriteResponse(220, "Ready to start TLS")

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConn := tls.Server(s.conn, tlsConfig)

	if err := tlsConn.Handshake(); err != nil {
		fmt.Printf("TLS握手失败: %v\n", err)
		return
	}

	s.conn = tlsConn
	s.reader = bufio.NewReader(tlsConn)
	s.writer = bufio.NewWriter(tlsConn)
	s.tlsEnabled = true

	// 重置会话状态
	s.resetTransaction()
}

// processMessage 处理接收到的邮件
func (s *SMTPSession) processMessage() {
	message := s.data.String()

	fmt.Printf("收到邮件:\n")
	fmt.Printf("  From: %s\n", s.from)
	fmt.Printf("  To: %v\n", s.to)
	fmt.Printf("  Size: %d bytes\n", len(message))

	// 保存邮件到文件（可选）
	s.saveMessage(message)

	// 转发邮件
	if s.config.RelayHost != "" && s.config.RelayTo != "" {
		if err := s.relayMessage(message); err != nil {
			fmt.Printf("转发邮件失败: %v\n", err)
			s.WriteResponse(451, "Local error in processing")
			return
		}
	}

	s.WriteResponse(250, "OK")
}

// saveMessage 保存邮件到文件
func (s *SMTPSession) saveMessage(message string) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("mail_%s_%s.eml", timestamp, s.from)
	filename = strings.ReplaceAll(filename, "<", "")
	filename = strings.ReplaceAll(filename, ">", "")
	filename = strings.ReplaceAll(filename, "@", "_at_")

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("保存邮件失败: %v\n", err)
		return
	}
	defer file.Close()

	file.WriteString(message)
	fmt.Printf("邮件已保存到: %s\n", filename)
}

// relayMessage 转发邮件到指定邮箱
func (s *SMTPSession) relayMessage(message string) error {
	fmt.Printf("正在转发邮件到: %s\n", s.config.RelayTo)

	// 添加转发头部
	forwardedMessage := s.addForwardHeaders(message)

	// 连接到转发SMTP服务器
	auth := smtp.PlainAuth("", s.config.RelayUser, s.config.RelayPassword, 
		strings.Split(s.config.RelayHost, ":")[0])

	err := smtp.SendMail(
		s.config.RelayHost,
		auth,
		s.config.RelayUser,
		[]string{s.config.RelayTo},
		[]byte(forwardedMessage),
	)

	if err != nil {
		return fmt.Errorf("SMTP发送失败: %v", err)
	}

	fmt.Printf("邮件转发成功\n")
	return nil
}

// addForwardHeaders 解析原邮件，提取关键信息后重组为易读的转发邮件
func (s *SMTPSession) addForwardHeaders(message string) string {
	msg, err := mail.ReadMessage(strings.NewReader(message))
	if err != nil {
		// 解析失败，原样返回
		return message
	}

	dec := new(mime.WordDecoder)
	decodeHeader := func(h string) string {
		v, err := dec.DecodeHeader(h)
		if err != nil {
			return h
		}
		return v
	}

	origFrom := decodeHeader(msg.Header.Get("From"))
	origTo := decodeHeader(msg.Header.Get("To"))
	origSubject := decodeHeader(msg.Header.Get("Subject"))
	origDate := msg.Header.Get("Date")
	if origDate == "" {
		origDate = time.Now().Format(time.RFC1123Z)
	}

	body := strings.TrimSpace(extractTextBody(
		msg.Header.Get("Content-Type"),
		msg.Header.Get("Content-Transfer-Encoding"),
		msg.Body,
	))

	var buf strings.Builder
	fmt.Fprintf(&buf, "From: %s\r\n", s.config.RelayUser)
	fmt.Fprintf(&buf, "To: %s\r\n", s.config.RelayTo)
	fmt.Fprintf(&buf, "Subject: %s\r\n", encodeMIMEHeader("[转发] "+origSubject))
	fmt.Fprintf(&buf, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	buf.WriteString("\r\n")
	buf.WriteString("---------- 转发邮件 ----------\r\n")
	fmt.Fprintf(&buf, "发件人: %s\r\n", origFrom)
	fmt.Fprintf(&buf, "日  期: %s\r\n", origDate)
	fmt.Fprintf(&buf, "主  题: %s\r\n", origSubject)
	fmt.Fprintf(&buf, "收件人: %s\r\n", origTo)
	buf.WriteString("\r\n")
	buf.WriteString(body)
	buf.WriteString("\r\n")
	return buf.String()
}

// encodeMIMEHeader 含非 ASCII 字符的头部使用 RFC 2047 编码
func encodeMIMEHeader(s string) string {
	for _, r := range s {
		if r > 127 {
			return "=?utf-8?B?" + base64.StdEncoding.EncodeToString([]byte(s)) + "?="
		}
	}
	return s
}

// extractTextBody 提取邮件正文的纯文本部分（解码 base64 / quoted-printable，展平 multipart）
func extractTextBody(contentType, encoding string, body io.Reader) string {
	if contentType == "" {
		return decodePart(body, encoding)
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return decodePart(body, encoding)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		return extractMultipart(body, params["boundary"])
	}
	return decodePart(body, encoding)
}

// extractMultipart 在 multipart 邮件中优先取 text/plain，否则降级 text/html 并剥离标签
func extractMultipart(body io.Reader, boundary string) string {
	if boundary == "" {
		data, _ := io.ReadAll(body)
		return string(data)
	}
	mr := multipart.NewReader(body, boundary)
	var plain, html string
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		ct := part.Header.Get("Content-Type")
		mediaType, params, _ := mime.ParseMediaType(ct)
		encoding := part.Header.Get("Content-Transfer-Encoding")

		if strings.HasPrefix(mediaType, "multipart/") {
			inner := extractMultipart(part, params["boundary"])
			if plain == "" {
				plain = inner
			}
			continue
		}

		decoded := decodePart(part, encoding)
		switch mediaType {
		case "text/plain":
			if plain == "" {
				plain = decoded
			}
		case "text/html":
			if html == "" {
				html = decoded
			}
		}
	}
	if plain != "" {
		return plain
	}
	return stripHTMLTags(html)
}

// decodePart 根据 Content-Transfer-Encoding 解码内容
func decodePart(body io.Reader, encoding string) string {
	var reader io.Reader = body
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "base64":
		reader = base64.NewDecoder(base64.StdEncoding, body)
	case "quoted-printable":
		reader = quotedprintable.NewReader(body)
	}
	data, _ := io.ReadAll(reader)
	return string(data)
}

// stripHTMLTags 简单去除 HTML 标签，仅在没有纯文本可用时作为兜底
func stripHTMLTags(html string) string {
	var b strings.Builder
	inTag := false
	for _, r := range html {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}
	return strings.TrimSpace(b.String())
}

// SMTPServer SMTP服务器
type SMTPServer struct {
	config *SMTPConfig
}

// NewSMTPServer 创建新的SMTP服务器
func NewSMTPServer(config *SMTPConfig) *SMTPServer {
	return &SMTPServer{config: config}
}

// Start 启动SMTP服务器
func (s *SMTPServer) Start() error {
	listener, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		return fmt.Errorf("监听失败: %v", err)
	}
	defer listener.Close()

	fmt.Printf("SMTP服务器已启动，监听地址: %s\n", s.config.ListenAddr)
	fmt.Printf("最大邮件大小: %d MB\n", s.config.MaxMessageSize/(1024*1024))
	
	if s.config.RelayHost != "" {
		fmt.Printf("转发配置:\n")
		fmt.Printf("  服务器: %s\n", s.config.RelayHost)
		fmt.Printf("  用户: %s\n", s.config.RelayUser)
		fmt.Printf("  目标邮箱: %s\n", s.config.RelayTo)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}

		go func() {
			session := NewSMTPSession(conn, s.config)
			session.Handle()
		}()
	}
}

func main() {
	// 配置
	config := &SMTPConfig{
		ListenAddr:     ListenAddr,
		TLSCertFile:    TLSCertFile,
		TLSKeyFile:     TLSKeyFile,
		MaxMessageSize: MaxMessageSize,
	}

	// 从环境变量读取转发配置
	config.RelayHost = os.Getenv("SMTP_RELAY_HOST")       // 如: smtp.gmail.com:587
	config.RelayUser = os.Getenv("SMTP_RELAY_USER")       // 如: your-email@gmail.com
	config.RelayPassword = os.Getenv("SMTP_RELAY_PASS")   // 邮箱密码或应用密码
	config.RelayTo = os.Getenv("SMTP_RELAY_TO")           // 目标转发邮箱

	// 从命令行参数覆盖
	if len(os.Args) > 1 {
		config.ListenAddr = os.Args[1]
	}

	server := NewSMTPServer(config)

	fmt.Println("Go SMTP Relay Server")
	fmt.Println("====================")
	fmt.Println("环境变量配置:")
	fmt.Println("  SMTP_RELAY_HOST    - 转发SMTP服务器地址")
	fmt.Println("  SMTP_RELAY_USER    - 转发邮箱用户名")
	fmt.Println("  SMTP_RELAY_PASS    - 转发邮箱密码")
	fmt.Println("  SMTP_RELAY_TO      - 目标转发邮箱地址")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  1. 设置环境变量后启动服务器")
	fmt.Println("  2. 配置邮件客户端发送邮件到此服务器")
	fmt.Println("  3. 邮件将被自动转发到指定邮箱")
	fmt.Println()

	if err := server.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}

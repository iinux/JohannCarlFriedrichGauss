package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "用法: eml-viewer <file.eml> [file.eml ...]")
		os.Exit(1)
	}

	for i, path := range os.Args[1:] {
		if i > 0 {
			fmt.Println(strings.Repeat("=", 60))
		}
		if err := viewFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "查看 %s 失败: %v\n", path, err)
		}
	}
}

func viewFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	msg, err := mail.ReadMessage(strings.NewReader(repairFoldedHeaders(string(raw))))
	if err != nil {
		return fmt.Errorf("解析邮件失败: %w", err)
	}

	dec := new(mime.WordDecoder)
	decodeHeader := func(h string) string {
		v, err := dec.DecodeHeader(h)
		if err != nil {
			return h
		}
		return v
	}

	fmt.Printf("文件: %s\n", path)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("发件人: %s\n", decodeHeader(msg.Header.Get("From")))
	fmt.Printf("收件人: %s\n", decodeHeader(msg.Header.Get("To")))
	if cc := msg.Header.Get("Cc"); cc != "" {
		fmt.Printf("抄  送: %s\n", decodeHeader(cc))
	}
	fmt.Printf("主  题: %s\n", decodeHeader(msg.Header.Get("Subject")))
	fmt.Printf("日  期: %s\n", msg.Header.Get("Date"))

	body, attachments := walk(
		msg.Header.Get("Content-Type"),
		msg.Header.Get("Content-Transfer-Encoding"),
		msg.Body,
	)

	if len(attachments) > 0 {
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("附件:")
		for _, a := range attachments {
			fmt.Printf("  - %s (%d 字节)\n", a.name, a.size)
		}
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("正文:")
	fmt.Println()
	fmt.Println(strings.TrimSpace(body))
	fmt.Println()
	return nil
}

type attachment struct {
	name string
	size int
}

// walk 递归处理邮件 body，返回首选正文文本和附件列表
func walk(contentType, encoding string, body io.Reader) (string, []attachment) {
	if contentType == "" {
		return decodePart(body, encoding), nil
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return decodePart(body, encoding), nil
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		return walkMultipart(body, params["boundary"])
	}
	return decodePart(body, encoding), nil
}

func walkMultipart(body io.Reader, boundary string) (string, []attachment) {
	if boundary == "" {
		data, _ := io.ReadAll(body)
		return string(data), nil
	}
	mr := multipart.NewReader(body, boundary)
	var plain, html string
	var attachments []attachment
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		ct := part.Header.Get("Content-Type")
		mediaType, params, _ := mime.ParseMediaType(ct)
		encoding := part.Header.Get("Content-Transfer-Encoding")
		disposition, dispParams, _ := mime.ParseMediaType(part.Header.Get("Content-Disposition"))

		// 附件分支
		if disposition == "attachment" || dispParams["filename"] != "" || params["name"] != "" {
			name := dispParams["filename"]
			if name == "" {
				name = params["name"]
			}
			dec := new(mime.WordDecoder)
			if decoded, err := dec.DecodeHeader(name); err == nil {
				name = decoded
			}
			data := decodePart(part, encoding)
			attachments = append(attachments, attachment{name: name, size: len(data)})
			continue
		}

		if strings.HasPrefix(mediaType, "multipart/") {
			innerBody, innerAtts := walkMultipart(part, params["boundary"])
			if plain == "" {
				plain = innerBody
			}
			attachments = append(attachments, innerAtts...)
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
		return plain, attachments
	}
	return stripHTMLTags(html), attachments
}

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

// repairFoldedHeaders 修复被 SMTP 服务器误删了前导空白的折叠头
// 头部区域内，非空行如果不是合法头部行（^name:）就视为续行，加 tab 前缀
// multipart 邮件中遇到 boundary 行后会重新进入头部块
func repairFoldedHeaders(raw string) string {
	lines := strings.Split(raw, "\n")
	inHeaders := true
	for i, line := range lines {
		trimmed := strings.TrimRight(line, "\r")

		if !inHeaders {
			// 遇到 multipart boundary 进入下一段头部块
			if strings.HasPrefix(trimmed, "--") {
				inHeaders = true
			}
			continue
		}

		if trimmed == "" {
			inHeaders = false
			continue
		}
		if trimmed[0] == ' ' || trimmed[0] == '\t' {
			continue
		}
		if looksLikeHeader(trimmed) {
			continue
		}
		// 视作被破坏的续行
		if strings.HasSuffix(line, "\r") {
			lines[i] = "\t" + trimmed + "\r"
		} else {
			lines[i] = "\t" + trimmed
		}
	}
	return strings.Join(lines, "\n")
}

// looksLikeHeader 判断是否是形如 "Header-Name: ..." 的头部行
func looksLikeHeader(line string) bool {
	for i, c := range line {
		if c == ':' {
			return i > 0
		}
		isAlnum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
		if !isAlnum && c != '-' && c != '_' {
			return false
		}
	}
	return false
}

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
	return b.String()
}

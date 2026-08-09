package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type EMLEmail struct {
	Headers  map[string]string
	From     string
	To       string
	Subject  string
	Body     string
	HTMLBody string
}

var softBreakRegex = regexp.MustCompile(`=\r?\n`)

func ParseEML(filename string) (*EMLEmail, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", filename, err)
	}

	email := &EMLEmail{Headers: make(map[string]string)}

	r := bufio.NewReader(strings.NewReader(string(data)))

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		idx := strings.Index(line, ":")
		if idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			email.Headers[key] = value
		}
	}

	if from := email.Headers["From"]; from != "" {
		addr, err := mail.ParseAddress(from)
		if err == nil {
			email.From = addr.Address
		} else {
			email.From = from
		}
	}

	if to := email.Headers["To"]; to != "" {
		addr, err := mail.ParseAddress(to)
		if err == nil {
			email.To = addr.Address
		} else {
			email.To = to
		}
	}

	email.Subject = email.Headers["Subject"]

	contentType := email.Headers["Content-Type"]
	ct, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		ct = "text/plain"
	}

	transferEncoding := email.Headers["Content-Transfer-Encoding"]

	body, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", filename, err)
	}

	bodyStr := string(body)

	switch {
	case strings.HasPrefix(ct, "text/html"):
		email.HTMLBody = decodeBody(bodyStr, transferEncoding)
	case strings.HasPrefix(ct, "text/plain"):
		email.Body = decodeBody(bodyStr, transferEncoding)
	case strings.HasPrefix(ct, "multipart/"):
		email.Body = parseMultipart(bodyStr, params["boundary"])
	}

	if email.Body == "" && email.HTMLBody == "" {
		email.Body = decodeBody(bodyStr, transferEncoding)
	}

	return email, nil
}

func decodeQuotedPrintable(s string) string {
	s = softBreakRegex.ReplaceAllString(s, "")
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '=' && i+2 < len(s) {
			hex := s[i+1 : i+3]
			if b, err := parseHexByte(hex); err == nil {
				result.WriteByte(b)
				i += 2
				continue
			}
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

func parseHexByte(s string) (byte, error) {
	var b byte
	for _, c := range s {
		var val byte
		switch {
		case c >= '0' && c <= '9':
			val = byte(c - '0')
		case c >= 'A' && c <= 'F':
			val = byte(c - 'A' + 10)
		case c >= 'a' && c <= 'f':
			val = byte(c - 'a' + 10)
		default:
			return 0, fmt.Errorf("invalid hex")
		}
		b = b<<4 | val
	}
	return b, nil
}

func decodeBody(body, encoding string) string {
	switch strings.ToLower(encoding) {
	case "quoted-printable":
		return decodeQuotedPrintable(body)
	case "base64":
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return body
		}
		return string(decoded)
	default:
		return body
	}
}

func parseMultipart(body, boundary string) string {
	if boundary == "" {
		return ""
	}
	parts := strings.Split(body, "--"+boundary)
	var textBody, htmlBody string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "--" {
			continue
		}

		headersEnd := strings.Index(part, "\r\n\r\n")
		if headersEnd == -1 {
			headersEnd = strings.Index(part, "\n\n")
		}
		if headersEnd == -1 {
			continue
		}

		headerStr := part[:headersEnd]
		partBody := part[headersEnd+4:]

		partHeaders := make(map[string]string)
		scanner := bufio.NewScanner(strings.NewReader(headerStr))
		for scanner.Scan() {
			line := scanner.Text()
			idx := strings.Index(line, ":")
			if idx > 0 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+1:])
				partHeaders[key] = value
			}
		}

		partCT := partHeaders["Content-Type"]
		partTE := partHeaders["Content-Transfer-Encoding"]

		ct, _, _ := mime.ParseMediaType(partCT)

		decodedBody := decodeBody(partBody, partTE)

		switch {
		case strings.HasPrefix(ct, "text/html"):
			htmlBody += decodedBody
		case strings.HasPrefix(ct, "text/plain"):
			textBody += decodedBody
		}
	}

	if htmlBody != "" {
		return htmlBody
	}
	return textBody
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: eml <eml_file> [eml_file2 ...]")
		fmt.Println("   or: eml <directory>")
		os.Exit(1)
	}

	args := os.Args[1:]

	var files []string
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		if info.IsDir() {
			entries, err := os.ReadDir(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading directory %s: %v\n", arg, err)
				continue
			}
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".eml") {
					files = append(files, filepath.Join(arg, entry.Name()))
				}
			}
		} else {
			files = append(files, arg)
		}
	}

	for _, file := range files {
		email, err := ParseEML(file)
		if err != nil {
			log.Printf("Error parsing %s: %v", file, err)
			continue
		}

		fmt.Printf("=== %s ===\n", filepath.Base(file))
		fmt.Printf("From:    %s\n", email.From)
		fmt.Printf("To:      %s\n", email.To)
		fmt.Printf("Subject: %s\n", email.Subject)

		body := email.HTMLBody
		if body == "" {
			body = email.Body
		}

		if len(body) > 200 {
			//body = body[:200] + "..."
		}
		body = strings.ReplaceAll(body, "\r\n", " ")
		body = strings.ReplaceAll(body, "\n", " ")
		fmt.Printf("Body:    %s\n", body)
		fmt.Println()
	}
}

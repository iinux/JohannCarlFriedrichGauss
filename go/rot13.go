package main

import (
    "io"
    "strings"
    "os"
)


type rot13Reader struct {
	r io.Reader
}

func (r *rot13Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err != nil || n == 0 {
		return
	}
	for i := 0; i < n; i++ {
		c := p[i]
		if c >= 'A' && c <= 'M' || c >= 'a' && c <= 'm' {
			c += 13
		} else {
			c -= 13
		}
		p[i] = c
	}
	return
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}

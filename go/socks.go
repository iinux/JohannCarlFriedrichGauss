package main

import (
"fmt"
"github.com/gamexg/proxyclient"
"io"
"io/ioutil"
)

func main() {
	p, err := proxyclient.NewProxyClient("socks5://127.0.0.1:1080")
	if err != nil {
		panic("创建代理客户端错误")
	}

	c, err := p.Dial("tcp", "www.google.com:80")
	if err != nil {
		panic("连接错误")
	}

	io.WriteString(c, "GET / HTTP/1.0\r\nHOST:www.google.com\r\n\r\n")
	b, err := ioutil.ReadAll(c)
	if err != nil {
		panic("读错误")
	}
	fmt.Print(string(b))
}

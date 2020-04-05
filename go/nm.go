package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const TOKEN = "hello"

func handleConn(c net.Conn) {
	fmt.Println(getTimeStr(), "start")
	buf := make([]byte, len(TOKEN))
	for true {
		n, err := c.Read(buf)
		if err == io.EOF {
			fmt.Println(getTimeStr(), c.RemoteAddr(), "end")
			break
		}
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		if n != len(TOKEN) {
			fmt.Println("read len error", n)
		}
		n, err = c.Write(buf)
		if err != nil {
			fmt.Println("write error", err)
			continue
		}
		if n != len(TOKEN) {
			fmt.Println("write len error", n)
		}
	}

	err := c.Close()
	if err != nil {
		fmt.Println("close error:", err)
		return
	}
}

func server(port int) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		go handleConn(c)
	}

}

func client(host string, port int, interval int) {
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), 10*time.Second)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	buf := make([]byte, len(TOKEN))
	var t time.Time
	var cost time.Duration
	var n int
	for true {
		time.Sleep(time.Duration(interval) * time.Second)

		t = time.Now()

		n, err = conn.Write([]byte(TOKEN))
		if err != nil {
			fmt.Println("write error:", err)
			break
		}
		if n != len(TOKEN) {
			fmt.Println("write len error:", n)
		}

		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		if n != len(TOKEN) {
			fmt.Println("read len error:", n)
		}

		cost = time.Since(t)
		fmt.Println(getTimeStr(), cost)
	}
}

func main() {
	mode := flag.String("m", "s", "mode server(s) or client(c)")
	host := flag.String("h", "127.0.0.1", "host")
	port := flag.Int("p", 8888, "port")
	internal := flag.Int("i", 1, "interval")
	flag.Parse()

	if *mode == "s" {
		server(*port)
	} else if *mode == "c" {
		client(*host, *port, *internal)
	} else {
		fmt.Println("error mode s or c")
	}
}

func getTimeStr() string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	return time.Now().In(cstSh).Format("2006-01-02 15:04:05")
}

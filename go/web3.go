package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"syscall"
	"runtime"
)

var (
	abort bool
)

const (
	SOCK = "/tmp/go.sock"
)

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body := "Hello World\n"
	// Try to keep the same amount of headers
	w.Header().Set("Server", "gophr")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	fmt.Fprint(w, body)
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	server := Server{}

	go func() {
		http.Handle("/", server)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		tcp, err := net.Listen("tcp", ":9001")
		if err != nil {
			log.Fatal(err)
		}
		fcgi.Serve(tcp, server)
	}()

	if runtime.GOOS != "windows" {
		go func() {
			unix, err := net.Listen("unix", SOCK)
			if err != nil {
				log.Fatal(err)
			}
			fcgi.Serve(unix, server)
		}()
	}

	<-sigchan

	if err := os.Remove(SOCK); err != nil {
		log.Fatal(err)
	}
}

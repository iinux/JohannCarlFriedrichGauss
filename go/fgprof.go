package main

import (
	_ "net/http/pprof"
	"github.com/felixge/fgprof"
	"net/http"
	"log"
	"time"
)

func main() {
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	for true {
		time.Sleep(1 * time.Second)

	}
	// <code to profile>
}

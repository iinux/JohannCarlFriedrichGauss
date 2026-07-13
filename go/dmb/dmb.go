// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/gxfont"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
)

// showMessage pops up a message box window. onClose, if non-nil, runs in
// addition to stopping the color animation when the window is closed.
func showMessage(driver gxui.Driver, theme gxui.Theme, font gxui.Font, title, content string, onClose func()) {
	window := theme.CreateWindow(50*len(content), 100, title)
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))

	label := theme.CreateLabel()
	label.SetFont(font)
	label.SetText(content)

	window.AddChild(label)

	ticker := time.NewTicker(time.Millisecond * 30)
	go func() {
		phase := float32(0)
		for _ = range ticker.C {
			c := gxui.Color{
				R: 0.75 + 0.25*math.Cosf((phase+0.000)*math.TwoPi),
				G: 0.75 + 0.25*math.Cosf((phase+0.333)*math.TwoPi),
				B: 0.75 + 0.25*math.Cosf((phase+0.666)*math.TwoPi),
				A: 0.50 + 0.50*math.Cosf(phase*10),
			}
			phase += 0.01
			driver.Call(func() {
				label.SetColor(c)
			})
		}
	}()

	window.OnClose(ticker.Stop)
	if onClose != nil {
		window.OnClose(onClose)
	}
}

// startServer opens a local HTTP port so a remote `curl` can pop up a
// message box, e.g. curl 'http://host:port/?t=title&c=hello+world'
// or curl --data 'hello world' http://host:port/
func startServer(driver gxui.Driver, theme gxui.Theme, font gxui.Font, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("t")
		if title == "" {
			title = "dmb"
		}
		content := r.URL.Query().Get("c")
		if content == "" && r.Method == http.MethodPost {
			body, _ := io.ReadAll(io.LimitReader(r.Body, 4096))
			content = string(body)
		}
		if content == "" {
			http.Error(w, "missing content: pass ?c=xxx or POST a body", http.StatusBadRequest)
			return
		}
		driver.Call(func() {
			showMessage(driver, theme, font, title, content, nil)
		})
		fmt.Fprintln(w, "ok")
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("dmb: listening on %s, e.g. curl 'http://<host>%s/?c=hello'", addr, addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Println("dmb: http server error:", err)
		}
	}()
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	font, err := driver.CreateFont(gxfont.Default, 75)
	if err != nil {
		panic(err)
	}


	if *port > 0 {
		startServer(driver, theme, font, *port)
	} else {
		showMessage(driver, theme, font, *title, *content, driver.Terminate)
	}
}

var title *string
var content *string
var port *int

func main() {
	// delay message box
	title = flag.String("t", "hello", "title")
	content = flag.String("c", "world", "content")
	var second = flag.Int("s", 0, "second")
	var minute = flag.Int("m", 0, "minute")
	var hour = flag.Int("h", 0, "hour")
	port = flag.Int("p", 0, "http port to listen on for remote curl requests (0 to disable), e.g. curl 'http://host:port/?c=hello'")

	flag.Parse()
	time.Sleep(time.Duration(*second)*time.Second + time.Duration(*minute)*time.Minute + time.Duration(*hour)*time.Hour)

	gl.StartDriver(appMain)
}

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"time"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/gxfont"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
)

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	font, err := driver.CreateFont(gxfont.Default, 75)
	if err != nil {
		panic(err)
	}

	window := theme.CreateWindow(50 * len(*content), 100, *title)
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))

	label := theme.CreateLabel()
	label.SetFont(font)
	label.SetText(*content)

	window.AddChild(label)

	ticker := time.NewTicker(time.Millisecond * 30)
	go func() {
		phase := float32(0)
		for _ = range ticker.C {
			c := gxui.Color{
				R: 0.75 + 0.25 * math.Cosf((phase + 0.000) * math.TwoPi),
				G: 0.75 + 0.25 * math.Cosf((phase + 0.333) * math.TwoPi),
				B: 0.75 + 0.25 * math.Cosf((phase + 0.666) * math.TwoPi),
				A: 0.50 + 0.50 * math.Cosf(phase * 10),
			}
			phase += 0.01
			driver.Call(func() {
				label.SetColor(c)
			})
		}
	}()

	window.OnClose(ticker.Stop)
	window.OnClose(driver.Terminate)
}

var title *string
var content *string

func main() {
	// delay message box
	title = flag.String("t", "hello", "title")
	content = flag.String("c", "world", "content")
	var second = flag.Int("s", 0, "second")
	var minute = flag.Int("m", 0, "minute")
	var hour = flag.Int("h", 0, "hour")

	flag.Parse()
	time.Sleep(time.Duration(*second) * time.Second + time.Duration(*minute) * time.Minute + time.Duration(*hour) * time.Hour)

	gl.StartDriver(appMain)
}

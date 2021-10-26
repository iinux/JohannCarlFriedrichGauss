package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {
	robotgo.ScrollMouse(10, "up")
	robotgo.Scroll(100, 200)
	robotgo.MilliSleep(100)
	robotgo.ScrollRelative(10, -100)

	robotgo.MouseSleep = 100
	robotgo.Move(10, 20)
	robotgo.MoveRelative(0, -10)
	robotgo.Drag(10, 10)

	robotgo.Click("left", true)
	robotgo.MoveSmooth(100, 200, 1.0, 100.0)
	robotgo.MouseToggle("up")
}
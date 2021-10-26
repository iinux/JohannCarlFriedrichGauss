package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
)

func main() {
	x, y := robotgo.GetMousePos()
	fmt.Println("pos: ", x, y)

	color := robotgo.GetPixelColor(100, 200)
	fmt.Println("color---- ", color)

	sx, sy := robotgo.GetScreenSize()
	fmt.Println("get screen size: ", sx, sy)

	bit := robotgo.CaptureScreen(10, 10, 30, 30)
	defer robotgo.FreeBitmap(bit)
	robotgo.SaveBitmap(bit, "test_1.png")

	img := robotgo.ToImage(bit)
	imgo.Save("test.png", img)
}
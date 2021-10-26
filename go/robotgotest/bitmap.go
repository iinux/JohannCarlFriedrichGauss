package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
	// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	defer robotgo.FreeBitmap(bitmap)

	fmt.Println("bitmap...", bitmap)
	img := robotgo.ToImage(bitmap)
	robotgo.SavePng(img, "test_1.png")

	bit2 := robotgo.ToCBitmap(robotgo.ImgToBitmap(img))
	fx, fy := robotgo.FindBitmap(bit2)
	fmt.Println("FindBitmap------ ", fx, fy)

	arr := robotgo.FindEveryBitmap(bit2)
	fmt.Println("Find every bitmap: ", arr)
	robotgo.SaveBitmap(bitmap, "test.png")

	fx, fy = robotgo.FindBitmap(bitmap)
	fmt.Println("FindBitmap------ ", fx, fy)

	robotgo.SaveBitmap(bitmap, "test.png")
}
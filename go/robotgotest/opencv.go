package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

func main() {
	opencv()
}

func opencv() {
	name := "test.png"
	name1 := "test_001.png"
	robotgo.SaveCapture(name1, 10, 10, 30, 30)
	robotgo.SaveCapture(name)

	fmt.Print("gcv find image: ")
	fmt.Println(gcv.FindImgFile(name1, name))
	fmt.Println(gcv.FindAllImgFile(name1, name))

	bit := robotgo.OpenBitmap(name1)
	defer robotgo.FindBitmap(bit)
	fmt.Print("find bitmap: ")
	fmt.Println(robotgo.FindBitmap(bit))

	// bit0 := robotgo.CaptureScreen()
	// img := robotgo.ToImage(bit0)
	// bit1 := robotgo.CaptureScreen(10, 10, 30, 30)
	// img1 := robotgo.ToImage(bit1)
	// defer robotgo.FreeBitmapArr(bit0, bit1)
	img := robotgo.CaptureImg()
	img1 := robotgo.CaptureImg(10, 10, 30, 30)

	fmt.Print("gcv find image: ")
	fmt.Println(gcv.FindImg(img1, img))
	fmt.Println(gcv.FindAllImg(img1, img))
}
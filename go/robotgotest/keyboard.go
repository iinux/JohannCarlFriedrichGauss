package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	robotgo.TypeStr("Hello World")
	robotgo.TypeStr("だんしゃり", 1.0)
	// robotgo.TypeStr("テストする")

	robotgo.TypeStr("Hi galaxy. こんにちは世界.")
	robotgo.Sleep(1)

	// ustr := uint32(robotgo.CharCodeAt("Test", 0))
	// robotgo.UnicodeType(ustr)

	robotgo.KeySleep = 100
	robotgo.KeyTap("enter")
	// robotgo.TypeStr("en")
	robotgo.KeyTap("i", "alt", "command")

	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)

	robotgo.MilliSleep(100)
	robotgo.KeyToggle("a", "down")

	robotgo.WriteAll("Test")
	text, err := robotgo.ReadAll()
	if err == nil {
		fmt.Println(text)
	}
}
package main

import (
	"time"
	"fmt"
	"syscall"
	"unsafe"
	"math/rand"
	"os"
	"strconv"
)

func abort(funcName string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcName, err))
}

var (
	user32, _      = syscall.LoadLibrary("user32.dll")
	messageBox, _  = syscall.GetProcAddress(user32, "MessageBoxW")
	keybd_event, _ = syscall.GetProcAddress(user32, "keybd_event")
)

const (
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_YESNOCANCEL       = 0x00000003
	MB_YESNO             = 0x00000004
	MB_RETRYCANCEL       = 0x00000005
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_ICONHAND          = 0x00000010
	MB_ICONQUESTION      = 0x00000020
	MB_ICONEXCLAMATION   = 0x00000030
	MB_ICONASTERISK      = 0x00000040
	MB_USERICON          = 0x00000080
	MB_ICONWARNING       = MB_ICONEXCLAMATION
	MB_ICONERROR         = MB_ICONHAND
	MB_ICONINFORMATION   = MB_ICONASTERISK
	MB_ICONSTOP          = MB_ICONHAND

	MB_DEFBUTTON1 = 0x00000000
	MB_DEFBUTTON2 = 0x00000100
	MB_DEFBUTTON3 = 0x00000200
	MB_DEFBUTTON4 = 0x00000300

	MB_SYSTEMMODAL = 0x00001000
)

func MessageBox(caption, text string, style uintptr) (result int) {
	var nargs uintptr = 4
	ret, _, callErr := syscall.Syscall9(uintptr(messageBox),
		nargs,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		style,
		0,
		0,
		0,
		0,
		0)
	if callErr != 0 {
		abort("Call MessageBox", callErr)
	}
	result = int(ret)
	return
}

func main() {
	wantMinuteStr := "0"
	var everySecond int64
	everySecond = 60
	if len(os.Args) >= 2 {
		wantMinuteStr = os.Args[1]
	}
	if len(os.Args) >= 3 {
		v, err := strconv.ParseInt(os.Args[2], 10, 0)
		if err == nil {
			everySecond = v
		} else {
			fmt.Println(err)
		}
	}
	tick := time.Tick(time.Duration(everySecond) * time.Second)

	for {
		<-tick
		now := time.Now()
		minuteStr := fmt.Sprintf("%d", now.Minute())
		if minuteStr == wantMinuteStr {
			MessageBox(now.Format("01-02 15:04"), random(1), MB_OK|MB_SYSTEMMODAL)
		}
	}
}

func random(category int) string {
	stringMap := map[int][]string{
		1: []string{
			"快速地敲击键盘, 让思维和效率飞奔起来",
			"不断学习",
			"学无止境",
			"不能浪费时间",
			"高效学习",
			"技术升华",
			"不断提升自己的技术层次",
		},
	}
	l := len(stringMap[category])
	r := rand.Intn(l)
	return stringMap[category][r]
}

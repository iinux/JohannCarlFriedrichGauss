package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	CF_UNICODETEXT = 13
)

var (
	user32 = syscall.MustLoadDLL("user32")
	openClipboard = user32.MustFindProc("OpenClipboard")
	closeClipboard = user32.MustFindProc("CloseClipboard")
	getClipboardData = user32.MustFindProc("GetClipboardData")
)

func main() {
	r, _, err := openClipboard.Call(0)
	if r == 0 {
		log.Fatalf("OpenClipboard failed: %v", err)
	}
	defer closeClipboard.Call()

	r, _, err = getClipboardData.Call(CF_UNICODETEXT)
	if r == 0 {
		log.Fatalf("GetClipboardData failed: %v", err)
	}
	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(r))[:])
	fmt.Printf("My clipboard has %q.\n", text)
}

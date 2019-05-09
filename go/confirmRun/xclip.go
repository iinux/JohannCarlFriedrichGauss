package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
	"net/http"
	"net/url"
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
)

const (
	CF_UNICODETEXT = 13
	SERVER = "http://192.168.188.221:8888"
)

func oldMain() {
	user32 := syscall.MustLoadDLL("user32")
	openClipboard := user32.MustFindProc("OpenClipboard")
	closeClipboard := user32.MustFindProc("CloseClipboard")
	getClipboardData := user32.MustFindProc("GetClipboardData")
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

func main() {
	args := os.Args
	var s walk.ClipboardService
	var clipText string

	if len(args) < 2 {
		fmt.Println("only support pull and push")
		return
	}

	if args[1] == "pull"  || args[1] == "l" {
		resp, err := http.PostForm(SERVER+"/get-clip",
			url.Values{"key": {"911"}})
		if err != nil {
			fmt.Println(err)
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				clipText = string(body)
				s.SetText(clipText)
			}
		}
	} else if args[1] == "push" || args[1] == "s" {
		clipText, err := s.Text()
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := http.PostForm(SERVER+"/set-clip",
			url.Values{"key": {"911"}, "text": {clipText}})
		if err != nil {
			fmt.Println(err)
		} else {
			defer resp.Body.Close()
			_, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			} else {
			}
		}
	} else {
		fmt.Println("only support pull and push")
	}
}

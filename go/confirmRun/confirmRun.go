package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"os"
	"os/exec"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	// "github.com/go-vgo/robotgo"
)

func abort(funcName string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcName, err))
}

var (
	user32, _     = syscall.LoadLibrary("user32.dll")
	messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")
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

func GetModuleHandle() (handle uintptr) {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(kernel32)
	getModuleHandle, _ := syscall.GetProcAddress(kernel32, "GetModuleHandleW")

	var nargs uintptr = 0
	if ret, _, callErr := syscall.Syscall(uintptr(getModuleHandle), nargs, 0, 0, 0); callErr != 0 {
		abort("Call GetModuleHandle", callErr)
	} else {
		handle = ret
	}
	return
}

func main() {
	defer syscall.FreeLibrary(user32)
	args := os.Args
	var text, program string
	program = args[0]
	text = "确定运行"

	var newArgs []string
	for k, s := range args {
		text += " " + s
		if k < 1 {
			continue
		}
		newArgs = append(newArgs, s)
	}

	text += " 吗?"
	program += "_real"

	userClick := MessageBox("Title", text, MB_YESNOCANCEL)
	fmt.Printf("Return: %d\n", userClick)

	if userClick == 6 {
		cmd := exec.Command(program, newArgs...)
		runCmd(cmd)
	} else if userClick == 7 {
		resp, err := http.PostForm("http://10.4.123.218:8888",
			url.Values{"key": {"911"}, "args": newArgs})
		if err != nil {
			fmt.Println(err)
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(body))
			}
			// robotgo.KeyTap("4", "command")
			// Replace the above line of code with the following code, the program size will be reduced by 1.6M
			// if use `go build -ldflags '-w -s -H=windowsgui'` will be reduced by 3.3M
			// https://docs.microsoft.com/en-us/windows/desktop/api/winuser/nf-winuser-keybd_event
			// https://docs.microsoft.com/zh-cn/windows/desktop/inputdev/virtual-key-codes
			// http://blog.51cto.com/kaixinbuliao/1348378

			syscall.Syscall(uintptr(keybd_event), 4, 91, 0, 0) // win key down
			syscall.Syscall(uintptr(keybd_event), 4, 52, 0, 0) // 4 key down
			syscall.Syscall(uintptr(keybd_event), 4, 52, 0, 2) // 4 key up
			syscall.Syscall(uintptr(keybd_event), 4, 91, 0, 2) // win key up
		}
		/*
		program = "C:\\Program Files\\Git\\usr\\bin\\ssh.exe"
		newArgs = []string{"shakespeare", "export DISPLAY=:1 ; /opt/google/chrome/chrome"}
		for k, s := range args {
			if k < 1 {
				continue
			}
			newArgs[1] = newArgs[1] + " '" + s + "'"
		}

		fmt.Println(program, newArgs)

		cmd := exec.Command(program, newArgs...)
		runCmd(cmd)
		*/
	}
}

func runCmd(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return string(opBytes)
}

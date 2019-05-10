package main

import (
	"flag"
	"fmt"
	"math/rand"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	messageBox              = user32.NewProc("MessageBoxW")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExA")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	keyboardHook            HHOOK

	keyboardPressCount = 0
	lastPressTime      = time.Now()

	everySecond                 *int
	wantMinute                  *int
	pressCountMinute            *int
	pressCountMinimum           *int
	messageBoxTypeProtectSecond *int
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

	WH_KEYBOARD_LL = 13
	WH_KEYBOARD    = 2
	WM_KEYDOWN     = 256
	WM_SYSKEYDOWN  = 260
	WM_KEYUP       = 257
	WM_SYSKEYUP    = 261
	WM_KEYFIRST    = 256
	WM_KEYLAST     = 264
	PM_NOREMOVE    = 0x000
	PM_REMOVE      = 0x001
	PM_NOYIELD     = 0x002
	WM_LBUTTONDOWN = 513
	WM_RBUTTONDOWN = 516
	NULL           = 0
)

func MessageBox(caption, text string, style uintptr) (result int) {
	for true {
		now := time.Now()
		if now.Sub(lastPressTime) < time.Duration(*messageBoxTypeProtectSecond)*time.Second {
			fmt.Println("becase you are in typing, waiting...")
			time.Sleep(time.Duration(*messageBoxTypeProtectSecond) * time.Second)
		} else {
			break
		}
	}

	messageBox.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		style)
	return
}

func main() {
	wantMinute = flag.Int("m", 0, "random string tip in which minute")
	everySecond = flag.Int("s", 60, "check random string tip every N second")
	pressCountMinute = flag.Int("pm", 1, "press count period(minute)")
	pressCountMinimum = flag.Int("pc", 10, "if press count low than this in period will alert")
	messageBoxTypeProtectSecond = flag.Int("mbp", 1, "if in typing message box delay second")
	flag.Parse()

	go PressCountStart()
	go PressCountTip()
	go RandomStrTip()

	for {
		time.Sleep(time.Hour)
	}
}

func RandomStrTip() {
	tick := time.Tick(time.Duration(*everySecond) * time.Second)
	for {
		<-tick
		now := time.Now()
		if now.Minute() == *wantMinute {
			MessageBox(now.Format("01-02 15:04"), RandomStr(1), MB_OK|MB_SYSTEMMODAL)
		}
	}
}

func RandomStr(category int) string {
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

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func PressCountStart() {
	// defer user32.Release()
	keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL,
		(HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
			if nCode == 0 && wparam == WM_KEYDOWN {
				lastPressTime = time.Now()
				if false {
					kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
					code := byte(kbdstruct.VkCode)
					fmt.Printf("key pressed: %q\n", code)
				}
				keyboardPressCount++
			}
			return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
		}), 0, 0)
	var msg MSG
	for GetMessage(&msg, 0, 0, 0) != 0 {

	}

	UnhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
}

func PressCountTip() {
	tick := time.Tick(time.Duration(*pressCountMinute) * time.Minute)
	for {
		<-tick

		fmt.Println(time.Now(), keyboardPressCount)
		if keyboardPressCount < *pressCountMinimum {
			pressStr := fmt.Sprintf("you have %d press in last %d minute(s)", keyboardPressCount, *pressCountMinute)
			MessageBox("Press Count Tip", pressStr, MB_OK|MB_SYSTEMMODAL)
		}

		keyboardPressCount = 0
	}
}

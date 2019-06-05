package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
	procOpenDesktop         = user32.NewProc("OpenDesktopW") // 用OpenDesktopA不行
	procCloseDesktop        = user32.NewProc("CloseDesktop")
	procSwitchDesktop       = user32.NewProc("SwitchDesktop")

	keyboardPressCount                        = 0
	lastPressTime                             = time.Now()
	mbSwitch                                  = true
	cpuUsage1min, cpuUsage5min, cpuUsage15min *CPUUsageLoad
	runTime                                   = time.Now()

	everySecond                 *int
	wantMinutesStr              *string
	pressCountMinute            *int
	pressCountMinimum           *int
	messageBoxTypeProtectSecond *int
	webListen                   *string
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

	DESKTOP_SWITCHDESKTOP = 0x0100
)

func MessageBox(caption, text string, style uintptr) (result int) {
	if !mbSwitch || GetWorkStationLocked() {
		return
	}

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
	wantMinutesStr = flag.String("m", "0|30", "random string tip in which minute")
	everySecond = flag.Int("s", 60, "check random string tip every N second")

	pressCountMinute = flag.Int("pm", 2, "press count period(minute)")
	pressCountMinimum = flag.Int("pc", 20, "if press count low than this in period will alert")

	webListen = flag.String("S", "", "web listen address and port")

	messageBoxTypeProtectSecond = flag.Int("mbp", 1, "if in typing message box delay second")
	flag.Parse()

	go CPUUsage()
	go PressCountStart()
	go PressCountTip()
	go RandomStrTip()

	if *webListen != "" {
		go Web()
	}

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "mb on" {
			mbSwitch = true
		} else if input == "mb off" {
			mbSwitch = false
		} else if input == "mb" {
			fmt.Println(mbSwitch)
		} else if input == "exit" || input == "quit" {
			os.Exit(0)
		} else if input == "uptime" {
			fmt.Println(runTime, time.Now().Sub(runTime))
		} else {
			fmt.Printf("%q is not a valid input\n", input)
		}
	}
}

func cpuLoad(w http.ResponseWriter, r *http.Request) {
	v, _ := mem.VirtualMemory()
	fmt.Fprintf(w, "%s -> Press: %3d, CPU: %6.2f %6.2f %6.2f Mem: %.0f%%\n",
		time.Now().Format("01-02 15:04:05"),
		keyboardPressCount,
		cpuUsage1min.Average,
		cpuUsage5min.Average,
		cpuUsage15min.Average,
		v.UsedPercent)
}

func Web() {
	http.HandleFunc("/", cpuLoad)
	http.HandleFunc("/cpuLoad", cpuLoad)
	fmt.Println("Listening on ", *webListen)
	err := http.ListenAndServe(*webListen, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

type CPUUsageLoad struct {
	Len     int
	Data    []float64
	Index   int
	Average float64
	Sum     float64
}

func (c *CPUUsageLoad) Append(value float64) {
	c.Sum -= c.Data[c.Index]
	c.Data[c.Index] = value
	c.Sum += c.Data[c.Index]
	c.Index = (c.Index + 1) % c.Len
	c.Average = c.Sum / float64(c.Len)
}

func NewCPUUsageLoad(len int, initValue float64) *CPUUsageLoad {
	c := CPUUsageLoad{}
	c.Len = len
	for i := 0; i < c.Len; i++ {
		c.Data = append(c.Data, initValue)
	}
	c.Sum = initValue * float64(c.Len)
	c.Average = initValue

	return &c
}

func CPUUsage() {
	percentage, _ := cpu.Percent(1, false)
	cpuUsage1min = NewCPUUsageLoad(60, percentage[0])
	cpuUsage5min = NewCPUUsageLoad(300, percentage[0])
	cpuUsage15min = NewCPUUsageLoad(900, percentage[0])

	go func() {

		tick := time.Tick(time.Second)
		for true {
			<-tick
			// cpu - get CPU number of cores and speed
			// cpuStat, err := cpu.Info()
			// fmt.Println(cpuStat, err)
			percentage, _ = cpu.Percent(0, false)
			cpuUsage1min.Append(percentage[0])
			cpuUsage5min.Append(percentage[0])
			cpuUsage15min.Append(percentage[0])
		}
	}()

	/*
		go func() {
			tick := time.Tick(time.Minute)
			for true {
				<-tick
				fmt.Println(percentage[0], _1min.Average, _5min.Average, _15min.Average)
			}
		}()
	*/
}

func GetWorkStationLocked() bool {
	r1, _, _ := procOpenDesktop.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Default"))),
		0, 0, DESKTOP_SWITCHDESKTOP)
	if r1 == 0 {
		fmt.Println("get desktop locked status error")
		return false
	}

	res, _, _ := procSwitchDesktop.Call(r1)
	// clean up
	procCloseDesktop.Call(r1)

	return res != 1
}

func RandomStrTip() {
	tick := time.Tick(time.Duration(*everySecond) * time.Second)
	for {
		<-tick
		now := time.Now()
		minuteStr := fmt.Sprintf("%d", now.Minute())
		wantMinuteStrArr := strings.Split(*wantMinutesStr, "|")
		for _, s := range wantMinuteStrArr {
			if minuteStr == s {
				MessageBox(now.Format("01-02 15:04"), RandomStr(1), MB_OK|MB_SYSTEMMODAL)
			}
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
			"莎士比亚也曾说过：“抛弃时间的人，时间也抛弃他。”人生短短数十载，千万不要“25岁死去，75岁埋葬”，荒废了一生。",
			"在清晨的时候，按下计算机的启动按钮，好像体验到了童年的一种感觉了。",
			"驱邪是件快乐的事，完成未完成的事也是快乐的",
			"不再害怕动手改造，反而要积极地去做，这样才会让生活更美好。",
			"年轻就像朝阳，容不得片刻怠慢",
			"为什么要选择今天堕落明天努力",
			"困了但没到点，可以选择做事。",
			"伙伴，我错了",
			"刻苦学习，增强体质，树立爱国主义、集体主义和社会主义思想，努力学习马克思列宁主义、毛泽东思想、邓小平理论，具有良好的思想品德，掌握较高的科学文化知识和专业技能。",
			"正规的任务不去做竟然浪费时间！",
			"我会为自己的慎独感到骄傲和自豪。",
			"与其荒诞，不如睡觉；与其犹豫，不如执行",
			"不要幻想打开潘多拉盒子之后能够合上",
			"今天家里第一次牵宽带了",
			"改过自新，重新做人",
			"自律，学习能力，创新能力，自信乐观，合作开放，责任感，执著追求，理性务实",
			"君子病无能焉 不病人之不己知也",
			"晚了能睡觉，起了能学习，幸福",
			"冲动是魔鬼，放纵是恶魔",
			"犹豫的时候直接把事情做了，不单快速地解决了问题，还顺带解决了犹豫的问题",
			"音乐的力量",
			"笔挻地为人和做事",
			"是非曲直",
			"白天能够学习，晚上能够睡觉其实是一件很幸福的事",
			"喷泉之所以漂亮是因为她有了压力；瀑布之所以壮观是因为她没有了退路；水之所以能穿石是因为永远在坚持。",
			"学习 做事 理智",
			"放纵自己得不到快乐，我不愿再放纵",
			"算法会了吗？",
			"格式塔心理学中有个名词叫“未完成事件”，它指的是尚未获得圆满解决或彻底弥合的既往情境，尤其是创伤或艰难情境，同时，也包含由此引发且未表达出来的情感，包括悔恨、愤怒、怨恨、痛苦、焦虑、悲伤、罪恶、遗弃感等。",
			"空谈误国、实干兴邦；时间就是金钱，效率就是生命",
			"解决问题方向很重要",
			"这种感觉是如此美妙，就像你刚刚掏出一支香烟，面前已是千百个打火机舞动。",
			"夫夷以近，则游者众；险以远，则至者少。而世之奇伟、瑰怪，非常之观，常在于险远，而人之所罕至焉，故非有志者不能至也。有志矣，不随以止也，然力不足者，亦不能至也。有志与力，而又不随以怠，至于幽暗昏惑而无物以相之，亦不能至也。然力足以至焉，于人为可讥，而在己为有悔；尽吾志也而不能至者，可以无悔矣，其孰能讥之乎？此余之所得也！王安石《游褒禅山记》",
			"高效率比花时间多更重要",
			"熬夜 罪已",
			"我要为我激情燃烧的岁月骄傲",
			"不赶不急得做事更稳",
			"痛快地做事感觉就是爽",
			"每天清晨，记得早起，努力追逐第一缕阳光，因为今天是我们余下的生命中最年轻的一天。",
			"要时刻清楚自己该干什么,不清楚到处乱撞肯定累.",
			"做个有思想的人",
			"要学会证明自己的实力，脱颖而出",
			"兰兰 眼睛 时间",
			"学好英语你可以和世界各地的程序员在Stack overflow，Reddit和Github进行交流，以码会友。",
			"有优秀的逻辑思维和抽象思维能力,有优秀的沟通技巧和情绪控制能力，可以无障碍地与开发、设计按照需求进行沟通，并保证项目顺利开展,有优秀的时间管理能力，保证项目开发进程可控,有良好的知识储备习惯,保持洞察一切（用户、产品、世界）的欲望，并具备独立思考精神……",
		},
	}
	l := len(stringMap[category])
	randMachine := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := randMachine.Intn(l)
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

		v, _ := mem.VirtualMemory()
		fmt.Printf("%s -> Press: %3d, CPU: %6.2f %6.2f %6.2f Mem: %.0f%%\n",
			time.Now().Format("01-02 15:04:05"),
			keyboardPressCount,
			cpuUsage1min.Average,
			cpuUsage5min.Average,
			cpuUsage15min.Average,
			v.UsedPercent)
		if keyboardPressCount < *pressCountMinimum {
			pressStr := fmt.Sprintf("you have %d press in last %d minute(s)", keyboardPressCount, *pressCountMinute)
			MessageBox("Press Count Tip", pressStr, MB_OK|MB_SYSTEMMODAL)
		}

		keyboardPressCount = 0
	}
}

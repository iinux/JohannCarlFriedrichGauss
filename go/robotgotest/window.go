package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	fpid, err := robotgo.FindIds("Google")
	if err == nil {
		fmt.Println("pids... ", fpid)

		if len(fpid) > 0 {
			robotgo.ActivePID(fpid[0])

			robotgo.Kill(fpid[0])
		}
	}

	robotgo.ActiveName("chrome")

	isExist, err := robotgo.PidExists(100)
	if err == nil && isExist {
		fmt.Println("pid exists is", isExist)

		robotgo.Kill(100)
	}

	abool := robotgo.ShowAlert("test", "robotgo")
	if abool {
		fmt.Println("ok@@@ ", "ok")
	}

	title := robotgo.GetTitle()
	fmt.Println("title@@@ ", title)
}
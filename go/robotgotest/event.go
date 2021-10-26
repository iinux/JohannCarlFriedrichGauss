package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func main() {
	add()
	low()
	event()
}

func add() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	robotgo.EventHook(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		robotgo.EventEnd()
	})

	fmt.Println("--- Please press w---")
	robotgo.EventHook(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("w")
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}

func event() {
	ok := robotgo.AddEvents("q", "ctrl", "shift")
	if ok {
		fmt.Println("add events...")
	}

	keve := robotgo.AddEvent("k")
	if keve {
		fmt.Println("you press... ", "k")
	}

	mleft := robotgo.AddEvent("mleft")
	if mleft {
		fmt.Println("you press... ", "mouse left button")
	}
}
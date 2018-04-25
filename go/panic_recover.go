package main

import (
	"os"
	"fmt"
)

var user = os.Getenv("USER")
func manualInit()  {
	if user == "" {
		panic("no value for $USER")
	}
	
}

func throwsPanic(f func()) (b bool) {
	defer func() {
		if x := recover(); x!= nil {
			b = true
		}
	}()
	f()
	return
}

func main()  {
	fmt.Println(throwsPanic(manualInit))
}

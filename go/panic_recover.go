package main

import (
	"os"
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
	throwsPanic(manualInit)
}

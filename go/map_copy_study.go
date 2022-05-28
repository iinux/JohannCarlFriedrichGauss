package main

import "fmt"

type A struct {
	I      int
	MapVar map[string]string
}

func main() {
	var a A
	a.MapVar = make(map[string]string)
	a.MapVar["hello"] = "world"
	a.I = 18

	ap := &a
	b := *ap

	b.MapVar["hello"] = "world2"
	b.I = 19

	c := &a.MapVar
	fmt.Println((*c)["hello"])
	fmt.Println(a.MapVar["hello"], a.I)
}

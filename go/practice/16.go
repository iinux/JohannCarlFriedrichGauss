package main

import "fmt"

type User struct {
	id   int
	name string
}

func main() {
	u := User{1, "Tom"}
	var vi, pi interface{} = u, &u

	// vi.(User).name = "Jack" // Error: cannot assign to vi.(User).name
	pi.(*User).name = "Jack"

	fmt.Printf("%v\n", vi.(User))
	fmt.Printf("%v\n", pi.(*User))
}

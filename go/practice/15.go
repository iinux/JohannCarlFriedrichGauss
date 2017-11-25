package main

import "fmt"

type User struct {
	id   int
	name string
}

func main() {
	u := User{1, "Tom"}
	var i interface{} = u

	u.id = 2
	u.name = "Jack"

	fmt.Printf("%v\n", u)
	fmt.Printf("%v\n", i.(User))
}

package main

import (
	"fmt"
)

type User struct {
	id   int
	name string
}

func (self *User) String() string {
	return fmt.Sprintf("%d, %s", self.id, self.name)
}

func main() {
	var o interface{} = &User{1, "Tom"}

	if i, ok := o.(fmt.Stringer); ok { // ok-idiom
		fmt.Println(i)
	}

	u := o.(*User)
	// u := o.(User) // panic: interface is *main.User, not main.User
	fmt.Println(u)
}

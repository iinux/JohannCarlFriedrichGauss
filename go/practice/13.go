package main

import "fmt"

type User struct {
	id   int
	name string
}

func (self *User) TestPointer() {
	fmt.Printf("TestPointer: %p, %v\n", self, self)
}

func (self User) TestValue() {
	fmt.Printf("TestValue: %p, %v\n", &self, self)
}

func main() {
	u := User{1, "Tom"}
	fmt.Printf("User: %p, %v\n", &u, u)

	mv := User.TestValue
	mv(u)

	mp := (*User).TestPointer
	mp(&u)

	mp2 := (*User).TestValue // *User 方法集包含 TestValue。
	mp2(&u)                  // 签名变为 func TestValue(self *User)。
} // 实际依然是 receiver value copy。

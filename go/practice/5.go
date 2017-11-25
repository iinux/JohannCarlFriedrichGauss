package main

type User struct {
	name string
	age  int
}

func main() {
	u1 := User{"Tom", 20}
	u2 := User{"Tom"} // Error: too few values in struct initializer
}

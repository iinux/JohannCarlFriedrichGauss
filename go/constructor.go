package main

type A struct {

}

func (*A) A()  {
	println("constructor")
}

func (*A) hello ()  {
	println("hello")
}

func main() {
	var a A
	a.hello()
}

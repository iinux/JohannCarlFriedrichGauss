package main

type Data struct{}

func (Data) TestValue()    {}
func (*Data) TestPointer() {}

func main() {
	var p *Data = nil
	p.TestPointer()

	(*Data)(nil).TestPointer() // method value
	(*Data).TestPointer(nil)   // method expression

	// p.TestValue() // invalid memory address or nil pointer dereference
	// (Data)(nil).TestValue() // cannot convert nil to type Data
	// Data.TestValue(nil) // cannot use nil as type Data in function argument
}

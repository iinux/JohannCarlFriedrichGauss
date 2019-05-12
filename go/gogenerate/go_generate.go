package main

import "fmt"

//go:generate echo hello
//go:generate go run go_generate.go
//go:generate  echo file=$GOFILE pkg=$GOPACKAGE
func main() {
	fmt.Println("main func")
}

//go:generate stringer -type=Pill
type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)

// refer https://www.jianshu.com/p/a866147021da

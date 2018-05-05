// go1.10-examples/spec/untypedconst.go
package main

import (
	"fmt"
)

var (
    s uint = 2
)

func main() {
    a := make([]int, 10)
    a[1.0<<s] = 4
	fmt.Println(a)
}

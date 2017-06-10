package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"C"
	"unsafe"
	"fmt"
)

func main() {
	fmt.Printf("hello, world\n")
	cstr := C.CString("hello world")
	C.puts(cstr)
	C.free(unsafe.Pointer(cstr))
}

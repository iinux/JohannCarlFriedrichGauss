package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import (
	"C"
	"unsafe"
)

func main() {
	cstr := C.CString("hello world")
	C.puts(cstr)
	C.free(unsafe.Pointer(cstr))
}

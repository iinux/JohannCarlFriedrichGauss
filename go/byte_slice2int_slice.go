package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

const SIZEOF_INT32 = 4 // bytes

func main() {
	// refer https://stackoverflow.com/questions/11924196/convert-between-slices-of-different-types#

	// The Right Way
	raw := []byte{1, 0, 0, 0, 1, 0, 0, 0}
	data := make([]int32, len(raw)/SIZEOF_INT32)
	for i := range data {
		// assuming little endian
		data[i] = int32(binary.LittleEndian.Uint32(raw[i*SIZEOF_INT32:(i+1)*SIZEOF_INT32]))
	}
	fmt.Println(data)

	// The Wrong Way
	// Get the slice header
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))

	// The length and capacity of the slice are different.
	header.Len /= SIZEOF_INT32
	header.Cap /= SIZEOF_INT32

	// Convert slice header to an []int32
	data1 := *(*[]int32)(unsafe.Pointer(&header))
	fmt.Println(data1)
}

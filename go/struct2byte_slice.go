package main

import (
	"fmt"
	"unsafe"
	"bytes"
	"encoding/binary"
	"math"
)

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type T struct {
	A int
	B string
}

/**
因为[]byte底层的数据结构为：
struct {
addr uintptr
len int
cap int
}
其中addr为数值的地址，len为当地数值的长度，cap为数值的容量。
转换的时候，需要定义一个和[]byte底层结构一致的struct（如例子中的SliceMock），
然后把结构体的地址赋给addr，结构体的大小赋给len和cap。最后将其转换为[]byte类型。
 */
func main() {
	var testStruct = &T{
		A:10,
		B:"abc",
	}
	Len := unsafe.Sizeof(*testStruct)
	testBytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(testStruct)),
		cap:  int(Len),
		len:  int(Len),
	}
	data := *(*[]byte)(unsafe.Pointer(testBytes))
	fmt.Println("Bytes:", data)

	// 1024 不能少
	pb := *(*[1024]byte)(unsafe.Pointer(testStruct))
	byteSlice := pb[:Len]
	fmt.Println("Bytes:", byteSlice)

	var ptestStruct *T = *(**T)(unsafe.Pointer(&data))
	fmt.Println("Struct: ", ptestStruct)

	pb2 := *(**T)(unsafe.Pointer(&byteSlice))
	fmt.Println("Struct: ", pb2)
}

/**
main2证明程序重新执行不可逆
 */

func main2()  {
	byteSlice := []byte{10, 0, 0, 0, 0, 0, 0, 0, 95, 102, 76, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0}
	fmt.Println("Bytes:", byteSlice)

	var ptestStruct *T = *(**T)(unsafe.Pointer(&byteSlice))
	fmt.Println("Struct: ", ptestStruct)

	pb2 := *(**T)(unsafe.Pointer(&byteSlice))
	fmt.Println("Struct: ", pb2)
}

func main3()  {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	err := binary.Write(buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
}

func main4()  {
	buf := new(bytes.Buffer)
	var testStruct = T{
		A:10,
		B:"abc",
	}
	// 结构体不可写入
	err := binary.Write(buf, binary.LittleEndian, &testStruct)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
}

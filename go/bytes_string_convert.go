package main

import "unsafe"

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}


func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // 获取s的起始地址开始后的两个 uintptr 指针
	h := [3]uintptr{x[0], x[1], x[1]}      // 构造三个指针数组
	return *(*[]byte)(unsafe.Pointer(&h))
}

func main()  {
	a := "hello"
	b := str2bytes(a)
	c := bytes2str(b)
	println(a,b,c)
}


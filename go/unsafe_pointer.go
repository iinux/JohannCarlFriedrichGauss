package main;

import (
    "fmt"
    "unsafe"
	"reflect"
)

func main() {
	i:= 10
	ip:=&i

	// error
	// var fp *float64 = (*float64)(ip)
    var fp *float64 = (*float64)(unsafe.Pointer(ip))

	*fp = *fp * 3
	fmt.Println(i)
	fmt.Println(fp)

    u:=new(user)
	fmt.Println(*u)

	pName:=(*string)(unsafe.Pointer(u))
	*pName="张三"

	pAge:=(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u))+unsafe.Offsetof(u.age)))
	*pAge = 20

    temp:=uintptr(unsafe.Pointer(u))+unsafe.Offsetof(u.age)
	pAge=(*int)(unsafe.Pointer(temp))
	*pAge = 21
    // 逻辑上看，以上代码不会有什么问题，但是这里会牵涉到GC，如果我们的这些临时变量被GC，那么导致的内存操作就错了，我们最终操作的，就不知道是哪块内存了，会引起莫名其妙的问题。

	fmt.Println(*u)

	// must use *u , u will error
	t:=reflect.TypeOf(*u)

	for i:=0;i<t.NumField();i++{
		sf:=t.Field(i)
		fmt.Println(sf.Tag)
		fmt.Println(sf.Tag.Get("haha"))
	}

	fmt.Println("======")

    fmt.Println(unsafe.Sizeof(true))
	fmt.Println(unsafe.Sizeof(int8(0)))
	fmt.Println(unsafe.Sizeof(int16(10)))
	fmt.Println(unsafe.Sizeof(int32(10000000)))
	fmt.Println(unsafe.Sizeof(int64(10000000000000)))

    // when GOARCH=386 error:  constant 10000000000000000 overflows int
	// fmt.Println(unsafe.Sizeof(int(10000000000000000)))
    // 对于和平台有关的int类型，这个要看平台是32位还是64位 GOARCH=386 4 GOARCH=amd64 8
	fmt.Println(unsafe.Sizeof(int(1000000000)))

	fmt.Println("======")

    var b bool
	var i8 int8
	var i16 int16
	var i64 int64
	var f32 float32
	var s string
	var m map[string]string
	var p *int32

	fmt.Println(unsafe.Alignof(b))
	fmt.Println(unsafe.Alignof(i8))
	fmt.Println(unsafe.Alignof(i16))
	fmt.Println(unsafe.Alignof(i64))
	fmt.Println(unsafe.Alignof(f32))
	fmt.Println(unsafe.Alignof(s))
	fmt.Println(unsafe.Alignof(m))
	fmt.Println(unsafe.Alignof(p))
    // 此外，获取对齐值还可以使用反射包的函数，也就是说：unsafe.Alignof(x)等价于reflect.TypeOf(x).Align()

	fmt.Println("======")

    var u1 user1
	fmt.Println(unsafe.Offsetof(u1.b))
	fmt.Println(unsafe.Offsetof(u1.i))
	fmt.Println(unsafe.Offsetof(u1.j))
	// 此外，unsafe.Offsetof(u1.i)等价于reflect.TypeOf(u1).Field(i).Offset

	fmt.Println("======")

	var u2 user2
	var u3 user3
	var u4 user4
	var u5 user5
	var u6 user6

	fmt.Println("u1 size is ",unsafe.Sizeof(u1))
	fmt.Println("u2 size is ",unsafe.Sizeof(u2))
	fmt.Println("u3 size is ",unsafe.Sizeof(u3))
	fmt.Println("u4 size is ",unsafe.Sizeof(u4))
	fmt.Println("u5 size is ",unsafe.Sizeof(u5))
	fmt.Println("u6 size is ",unsafe.Sizeof(u6))
}

type user struct {
	name string `haha:"name"`
	age int `haha:"age"`
}

type user1 struct {
	b byte
	i int32
	j int64
}

type user2 struct {
	b byte
	j int64
	i int32
}

type user3 struct {
	i int32
	b byte
	j int64
}

type user4 struct {
	i int32
	j int64
	b byte
}

type user5 struct {
	j int64
	b byte
	i int32
}

type user6 struct {
	j int64
	i int32
	b byte
}

// refer https://www.flysnow.org/2017/07/06/go-in-action-unsafe-pointer.html
// refer https://www.flysnow.org/2017/07/02/go-in-action-unsafe-memory-layout.html
// refer https://www.flysnow.org/2017/06/25/go-in-action-struct-tag.html


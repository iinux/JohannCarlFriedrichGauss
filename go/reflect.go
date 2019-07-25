package main

import (
	"fmt"
	"reflect"
	"strconv"
	"runtime"
)

func hello() {
	fmt.Println("Hello world!")
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func main() {
    u:= MyType{20, "张三"}
	t:=reflect.TypeOf(u)
	fmt.Println(t)
	v := reflect.ValueOf(u)
	fmt.Println(v)
	t1:=v.Type()
	fmt.Println(t1)
	fmt.Printf("%T\n",u)
	fmt.Printf("%v\n",u)

    for i:=0;i<t.NumField();i++ {
		fmt.Println(t.Field(i).Name)
	}

	for i:=0;i<t.NumMethod() ;i++  {
		fmt.Println(t.Method(i).Name)
	}

    // refer https://www.flysnow.org/2017/06/13/go-in-action-go-reflect.html
    fmt.Println("======")

	fmt.Println(GetFunctionName(hello))
	hl := hello
	fv := reflect.ValueOf(hl)
	fmt.Println("fv is reflect.Func ?", fv.Kind() == reflect.Func)
	fv.Call(nil)
	///
	fv2 := reflect.ValueOf(prints)
	params := make([]reflect.Value, 1)                 //参数
	params[0] = reflect.ValueOf(20)                    //参数设置为20
	rs := fv2.Call(params)                             //rs作为结果接受函数的返回值
	fmt.Println("result:", rs[0].Interface().(string)) //当然也可以直接是rs[0].Interface()
	///
	myType := &MyType{22, "wowzai"}
	//fmt.Println(myType)     //就是检查一下myType对象内容
	//println("---------------")
	mtV := reflect.ValueOf(&myType).Elem()
	fmt.Println("Before:", mtV.MethodByName("String").Call(nil)[0])
	params3 := make([]reflect.Value, 1)
	params3[0] = reflect.ValueOf(18)
	mtV.MethodByName("SetI").Call(params3)
	params3[0] = reflect.ValueOf("reflection test")
	mtV.MethodByName("SetName").Call(params3)
	fmt.Println("After:", mtV.MethodByName("String").Call(nil)[0])
	///
	fmt.Println("Before:", mtV.Method(2).Call(nil)[0])
	params[0] = reflect.ValueOf(18)
	mtV.Method(0).Call(params)
	params[0] = reflect.ValueOf("reflection test")
	mtV.Method(1).Call(params)
	fmt.Println("After:", mtV.Method(2).Call(nil)[0])
}

func prints(i int) string {
	fmt.Println("i =", i)
	return strconv.Itoa(i)
}

type MyType struct {
	i    int
	name string
}

func (mt *MyType) SetI(i int) {
	mt.i = i
}

func (mt *MyType) SetName(name string) {
	mt.name = name
}

func (mt *MyType) String() string {
	return fmt.Sprintf("%p", mt) + "--name:" + mt.name + " i:" + strconv.Itoa(mt.i)
}

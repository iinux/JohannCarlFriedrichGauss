package main

import "fmt"

type Power struct{
        age int
        high int
        name string
}

//指针类型
func (this *Power) String() string {
        return fmt.Sprintf("age:%d, high:%d, name:%s", this.age, this.high, this.name)
}


func main() {

        var i Power = Power{age: 10, high: 178, name: "NewMan"} //非指针

        fmt.Printf("%s\n", i)
        fmt.Println(i)
        fmt.Printf("%v", i)
}

// 其他三种情况都是OK的
// refer https://blog.csdn.net/lanyang123456/article/details/78178183

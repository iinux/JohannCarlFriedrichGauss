package main

import (
    "encoding/gob"
    "fmt"
    "os"
    "runtime"
)

type User struct {
    Id   int
    Name string
}

func (this *User) Say() string {
    return this.Name + ` hello world ! `
}

func main(){
    var filePath string
    if runtime.GOOS == "windows" {
		filePath = "d:/gob"
	} else {
		filePath = "mygo/gob"
	}
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Println(err)
    }
    user := User{Id: 1, Name: "Mike"}
    user2 := User{Id: 3, Name: "Jack"}
    u := []User{user, user2}
    enc := gob.NewEncoder(file)
    err2 := enc.Encode(u)
    fmt.Println(err2)
}

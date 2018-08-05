package main;

import "fmt"
import "github.com/kevinburke/nacl"
import "github.com/kevinburke/nacl/secretbox"
import "encoding/base64"

func main() {
    key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
    if err != nil {
        panic(err)
    }
    encrypted := secretbox.EasySeal([]byte("hello world"), key)
    fmt.Println(base64.StdEncoding.EncodeToString(encrypted))
}

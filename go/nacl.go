package main;

import "fmt"
import "github.com/kevinburke/nacl"
import "github.com/kevinburke/nacl/secretbox"
import "encoding/base64"

// http://nacl.cr.yp.to/
func main() {
    key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
    if err != nil {
        panic(err)
    }
    encrypted := secretbox.EasySeal([]byte("hello world"), key)
    fmt.Println(base64.StdEncoding.EncodeToString(encrypted))
    decrypted, err := secretbox.EasyOpen(encrypted, key)
    if err == nil {
        fmt.Println(string(decrypted))
    } else {
        fmt.Println(err)
    }
}

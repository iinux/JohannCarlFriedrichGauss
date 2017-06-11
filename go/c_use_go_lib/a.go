package main
import "C"
import "fmt"
//export Hello
func Hello() string {
    fmt.Println("Hello");
    return "Hello"
}
func main() {
}

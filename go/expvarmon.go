package main

import (
	_ "expvar"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello, world!")
	http.ListenAndServe(":1234", nil)
}

package main

import (
	"github.com/bcicen/grmon"
	"net/http"
)

func main()  {
	grmon.Start()
	http.ListenAndServe(":1234", nil)

}

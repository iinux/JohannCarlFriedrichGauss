package main

/*
const char* build_time(void)
{
        static const char* psz_build_time = "["__DATE__ "  " __TIME__ "]";
            return psz_build_time;
}
*/
import "C"

import (
	"fmt"
)

var (
	buildTime = C.GoString(C.build_time())
)

func BuildTime() string {
	return buildTime
}

func main() {
	fmt.Printf("Build time is: %s\n", BuildTime())
}

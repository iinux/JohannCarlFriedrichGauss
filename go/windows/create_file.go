package main

import (
	"fmt"
	"github.com/samhsu427/winlabs/gowin32"
	"os"
)

func main()  {
	if len(os.Args) < 2 {
		fmt.Println("not enough args")
		os.Exit(1)
	}
	fd, err := gowin32.OpenWindowsFile(os.Args[1], true, gowin32.FileShareRead, gowin32.FileOpenExisting, 0, gowin32.FileFlagPOSIXSemantics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := make([]byte, 6)
	n, err := fd.Read(b)
	if n == 0 {
		fmt.Println("read none")
		os.Exit(0)
	}

	fmt.Println(string(b))
}

package main

import (
	"github.com/sqweek/dialog"
	"fmt"
)

func main()  {
	ok := dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
	fmt.Println(ok)
}

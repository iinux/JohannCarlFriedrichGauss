package main

import (
	"fmt"
	"time"
	wapi "github.com/iamacarpet/go-win64api"
)

func main(){
	users, err := wapi.ListLocalUsers()
	if err != nil {
		fmt.Printf("Error fetching user list, %s.\r\n", err.Error())
		return
	}

	for true {
		fmt.Println(time.Now())
		for _, u := range users {
			fmt.Printf("%s (%s)\r\n", u.Username, u.FullName)
			fmt.Printf("\tIs Locked:                    %t\r\n", u.IsLocked)
			/*
		fmt.Printf("%s (%s)\r\n", u.Username, u.FullName)
		fmt.Printf("\tIs Enabled:                   %t\r\n", u.IsEnabled)
		fmt.Printf("\tIs Locked:                    %t\r\n", u.IsLocked)
		fmt.Printf("\tIs Admin:                     %t\r\n", u.IsAdmin)
		fmt.Printf("\tPassword Never Expires:       %t\r\n", u.PasswordNeverExpires)
		fmt.Printf("\tUser can't change password:   %t\r\n", u.NoChangePassword)
		fmt.Printf("\tPassword Age:                 %.0f days\r\n", (u.PasswordAge.Hours()/24))
		fmt.Printf("\tLast Logon Time:              %s\r\n", u.LastLogon.Format(time.RFC850))
		fmt.Printf("\tBad Password Count:           %d\r\n", u.BadPasswordCount)
		fmt.Printf("\tNumber Of Logons:             %d\r\n", u.NumberOfLogons)
		*/
		}
		time.Sleep(time.Minute)
	}
}

package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"log"
	"net/http"
	"strings"
	"encoding/json"
)

type Login struct {
	UserId int
	Username string
	Password string
}

type SendSms struct {
	ReceiptUid int
	Content string
}

type OnlineUser struct {
	Uid int
	MessageChannel chan string
}

var (
	OnlineUsers []OnlineUser
)

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		var index = strings.Index(reply, ":")
		var action = reply[:index]
		var rawData = reply[index + 1:]

		switch action {
		case "login":
			var data Login
			json.Unmarshal([]byte(rawData), &data)

			var user OnlineUser
			user.Uid = data.UserId
			user.MessageChannel = make(chan string)

			OnlineUsers = append(OnlineUsers, user)

			fmt.Println(data)
			fmt.Println(OnlineUsers)

			go func(){
				for {
					var msg string
					msg = <- user.MessageChannel
					if err = websocket.Message.Send(ws, msg); err != nil {
						fmt.Println("Can't send")
						break
					}
				}
			}()
		case "logout":
			fmt.Println("logout")
			break
		case "send":
			var data SendSms
			json.Unmarshal([]byte(rawData), &data)

			for _, onlineUser := range OnlineUsers {
				if onlineUser.Uid == data.ReceiptUid {
					onlineUser.MessageChannel <- data.Content
				}
			}
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
)

var input string

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:7978", Path: "/serveWs"}
	fmt.Println("connecting to %s\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("连接失败：", err.Error())
		return
	}
	defer c.Close()
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Client接收消息失败:\n", err.Error())
				return
			}
			fmt.Printf("Client接收消息:%s\n", msg)
		}
	}()

	go func() {
		for {
			n, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println(n, "--", err.Error())
			} else {
				msg := Msg{
					SelfId:  "1",
					PeerId:  "2",
					Cmd:     "sendMsg",
					MsgType: "text",
					Text:    input,
				}
				b, err := json.Marshal(&msg)
				if err != nil {
					fmt.Println("marshal fail...", err.Error())
					return
				}
				err = c.WriteMessage(websocket.TextMessage, b)
				if err != nil {
					fmt.Println("client", err.Error())
					return
				}
			}
		}
	}()
	//ticker := time.NewTicker(2 * time.Second)
	//defer ticker.Stop()
	for {
		select {
		//case <-ticker.C:
		//	msg := Msg{
		//		SelfId:"1",
		//		PeerId:"2",
		//		Cmd:     "login",
		//		MsgType: "text",
		//		Text:    "this is send text",
		//	}
		//	b, err := json.Marshal(&msg)
		//	if err != nil {
		//		fmt.Println("marshal fail...", err.Error())
		//		return
		//	}
		//	fmt.Println(msg)
		//	//s := "{\"Cmd\":\"sendMsg\",\"msgType\":\"text\",\"text\":\"this is send text\"}"
		//	err = c.WriteMessage(websocket.TextMessage, b)
		//	if err != nil {
		//		fmt.Println("Client写入数据失败：", err.Error())
		//		return
		//	}
		case <-interrupt:
			fmt.Println("客户端interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("客户端写入关闭数据失败：", err.Error())
				return
			}
		}
	}
}

type Msg struct {
	SelfId  string
	PeerId  string
	Cmd     string
	Text    string
	MsgType string
}

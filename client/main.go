package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

func main() {
	u := url.URL{
		Scheme: "ws",
		Host: "127.0.0.1:7978",
		Path: "/ws",
	}
	fmt.Println(u.String())
	c,_,err := websocket.DefaultDialer.Dial(u.String(),nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = c.WriteMessage(websocket.TextMessage,[]byte("测试数据..."))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

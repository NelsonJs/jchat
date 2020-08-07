package socketservice

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"jchat/logs"
	"time"
)

const (
	Android = 1
	IOS     = 2
	WEB     = 3
)

type Client struct {
	UserId    string
	LoginTime int64
	Socket    *websocket.Conn
	Addr      string //客户端地址
	AppId     int8   //登陆平台的id android/ios/web
	msg       chan []byte
}

func NewClient(loginTime int64, conn *websocket.Conn) *Client {
	return &Client{
		LoginTime: loginTime,
		Socket:    conn,
		msg:       make(chan []byte, 10),
	}
}

func (client *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger().Error(fmt.Sprint(r))
		}
	}()
	for {
		msgType, msg, err := client.Socket.ReadMessage()
		m := Msg{}
		err = json.Unmarshal(msg, &m)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(m)
		if err != nil {
			logs.Logger().Error(err.Error(), zap.String("read msg", "meet error....."),
				zap.Duration("time", time.Second))
			return
		}
		client.Process(msgType, msg)
	}
}

func (client *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger().Error(fmt.Sprint(r))
		}
	}()
	for {
		select {
		case b, ok := <-client.msg:
			fmt.Println("write..", ok, string(b))
			if !ok {
				return
			}
			client.Socket.WriteMessage(websocket.TextMessage, b)
		}
	}

}

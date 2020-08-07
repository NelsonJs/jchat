package socketservice

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"jchat/logs"
	"jchat/models"
)

const (
	TEXT  = "text"
	IMAGE = "img"
)

const (
	Login = "login"
	Send  = "sendMsg"
)

var cmdHandlers map[string]MsgHandler

type Msg struct {
	Cmd     string `json:"cmd"`
	SelfId  string
	PeerId  string
	Time    int64
	MsgId   string
	MsgType string
	Text    string
	Data    []byte
}

type MsgHandler func(client *Client, msg *Msg)

func putCmdFunc(cmd string, handler MsgHandler) {
	if cmdHandlers == nil {
		cmdHandlers = make(map[string]MsgHandler)
	}
	cmdHandlers[cmd] = handler
}

//注册函数
func init() {
	putCmdFunc(Send, sendMsg)
	putCmdFunc(Login, login)
}

func getCmdFunc(cmd string) (MsgHandler, bool) {
	if cmdHandlers == nil {
		return nil, false
	}
	v, ok := cmdHandlers[cmd]
	return v, ok
}

func (client *Client) Process(socketTransferType int, msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger().Error(fmt.Sprint(r))
		}
	}()

	msgModel := Msg{}
	err := json.Unmarshal(msg, &msgModel)
	if err != nil {
		logs.Logger().Error(err.Error())
		res := models.Response{Code: models.StatusMsgTypeError, Msg: models.StatusCodeTxt(models.StatusMsgTypeError)}
		b, _ := json.Marshal(&res)
		client.msg <- b
		return
	}
	//每个人发送消息之前，都需要验证 己方及对方 是否注册
	if msgModel.Cmd != Login {
		if !isLogin(msgModel.SelfId) {
			b, err := models.GetResponseByte(models.StatusNotLoginError, models.StatusCodeTxt(models.StatusNotLoginError))
			if err != nil {
				logs.Logger().Error(err.Error())
				return
			}
			client.msg <- b
			return
		}
		if !isLogin(msgModel.PeerId) {
			b, err := models.GetResponseByte(models.StatusReceiveUserNotLoginError, models.StatusCodeTxt(models.StatusReceiveUserNotLoginError))
			if err != nil {
				logs.Logger().Error(err.Error())
				return
			}
			client.msg <- b
			return
		}
	}
	if handler, ok := getCmdFunc(msgModel.Cmd); ok {
		handler(client, &msgModel)
	}
	fmt.Println("msgType--->", socketTransferType)
	if socketTransferType == websocket.TextMessage {

	} else if socketTransferType == websocket.BinaryMessage {

	} else if socketTransferType == websocket.PingMessage {

	} else if socketTransferType == websocket.PongMessage {

	} else if socketTransferType == websocket.CloseMessage {

	}
}

func sendMsg(client *Client, msg *Msg) {
	if msg.MsgType == TEXT {
		fmt.Println(msg.Text)
		c, ok := getClient(msg.PeerId)
		if ok {
			c.msg <- []byte(msg.Text)
		} else {
			b, err := models.GetResponseByte(models.StatusReceiveUserNotLoginError, models.StatusCodeTxt(models.StatusReceiveUserNotLoginError))
			if err != nil {
				logs.Logger().Error(err.Error())
				return
			}
			client.msg <- b
		}

	}
}

func login(client *Client, msg *Msg) {
	if msg != nil && client != nil && msg.SelfId != "" {
		client.UserId = msg.SelfId
		clientManager.Login <- client
	} else {
		b, err := models.GetResponseByte(models.StatusParametersError, models.StatusCodeTxt(models.StatusParametersError))
		if err != nil {
			logs.Logger().Error(err.Error())
			return
		}
		client.msg <- b
	}
}

package socketservice

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"jchat/logs"
	"jchat/models"
)
type MsgType int64
const (
	TEXT MsgType = 1
	IMAGE MsgType = 2
)

type CmdType int64

const (
	Login CmdType = 1
	Send CmdType = 2
)

var cmdHandlers map[CmdType]MsgHandler

type Msg struct {
	cmd CmdType
	selfId string
	peerId string
	time int64
	msgId string
	msgType MsgType
	text string
	data []byte
}

type MsgHandler func(client *Client,msg *Msg)

func putCmdFunc(cmd CmdType,handler MsgHandler)  {
	if cmdHandlers == nil {
		cmdHandlers = make(map[CmdType]MsgHandler)
	}
	cmdHandlers[cmd] = handler
}

func init() {
	putCmdFunc(Send,SendMsg)
}

func getCmdFunc(cmd CmdType) (MsgHandler,bool) {
	if cmdHandlers == nil {
		return nil,false
	}
	v,ok := cmdHandlers[cmd]
	return v,ok
}

func (client *Client) Process(socketTransferType int,msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			logs.Logger().Error(fmt.Sprint(r))
		}
	}()

	msgModel := Msg{}
	err := json.Unmarshal(msg,&msgModel)
	if err != nil {
		logs.Logger().Error(err.Error())
		res := models.Response{Code: models.StatusMsgTypeError,Msg: models.StatusCodeTxt(models.StatusMsgTypeError)}
		b,_ := json.Marshal(&res)
		client.msg <- b
		return
	}

	if handler,ok := getCmdFunc(msgModel.cmd); ok {
		handler(client,&msgModel)
	}

	if socketTransferType == websocket.TextMessage {

	} else if socketTransferType == websocket.BinaryMessage {

	} else if socketTransferType == websocket.PingMessage {

	} else if socketTransferType == websocket.PongMessage {

	} else if socketTransferType == websocket.CloseMessage {

	}
}


func SendMsg(client *Client,msg *Msg) {
	if msg.msgType == TEXT {
		client.msg <- []byte(msg.text)
	}

}

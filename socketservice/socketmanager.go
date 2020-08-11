package socketservice

import (
	"fmt"
	"github.com/gorilla/websocket"
	"jchat/logs"
	"jchat/models"
	"sync"
	"time"
)

type ClientManager struct {
	Clients       map[string]*Client //在线用户
	Login         chan *Client       //登录
	Logout        chan *Client       //下线
	Connection    chan *Client       //连接
	DisConnection chan *Client       //断开连接
}

var clientManager *ClientManager

func init() {
	clientManager = &ClientManager{
		Clients:       make(map[string]*Client),
		Login:         make(chan *Client, 100),
		Logout:        make(chan *Client, 100),
		Connection:    make(chan *Client, 100),
		DisConnection: make(chan *Client, 100),
	}
}

var mutex sync.Mutex

//登录
func (manager *ClientManager) register(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	if client.UserId != "" {
		manager.Clients[client.UserId] = client
		b, err := models.GetResponseByte(models.StatusLoginSuccess, models.StatusCodeTxt(models.StatusLoginSuccess))
		if err != nil {
			logs.Logger().Error(err.Error())
			return
		}
		client.msg <- b
		// 检测无响应的连接
		//每5秒发送一个ping消息，收到后回复一个pong消息
		if client.UserId == "2" {
			client.Socket.SetPingHandler(func(appData string) error {
				fmt.Println("ping handler--->", appData)
				client.Socket.WriteMessage(websocket.PongMessage, []byte("pong message.."))
				client.HeartBeatTime = time.Now().Unix()
				return nil
			})
		}

	}
}

func checkNoResponseConn() {
	fmt.Println("checking..")
	for _, v := range clientManager.Clients {
		fmt.Println("checking..", v.UserId)
		if v.UserId == "2" && v.clientExpired() {
			fmt.Println("del--->", v.UserId)
			delete(clientManager.Clients, v.UserId)
		}
	}
}

//下线
func (manager *ClientManager) logout(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	if client.UserId != "" {
		delete(manager.Clients, client.UserId)
		b, err := models.GetResponseByte(models.StatusLoginError, models.StatusCodeTxt(models.StatusLoginError))
		if err != nil {
			logs.Logger().Error(err.Error())
			return
		}
		client.msg <- b
	}
}

// 是否注册
func isLogin(userId string) bool {
	_, exists := clientManager.Clients[userId]
	return exists
}

// 获取登录过的用户
func getClient(userId string) (*Client, bool) {
	c, ok := clientManager.Clients[userId]
	return c, ok
}

var pingWait = 8 * time.Second
var checkWait = 10 * time.Second

func StartListen() {
	var t = time.NewTicker(checkWait)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			checkNoResponseConn()
		case c, ok := <-clientManager.Login:
			if ok {
				clientManager.register(c)
			}
		case c, ok := <-clientManager.Logout:
			if ok {
				clientManager.register(c)
			}
		}
	}
}

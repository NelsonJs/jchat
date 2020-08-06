package socketservice

import "sync"

type ClientManager struct {
	Clients map[string]*Client //在线用户
	Login chan *Client //登录
	Logout chan *Client //下线
	Connection chan *Client //连接
	DisConnection chan *Client //断开连接
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients: make(map[string]*Client),
		Login: make(chan *Client,100),
		Logout: make(chan *Client,100),
		Connection: make(chan *Client,100),
		DisConnection: make(chan *Client,100),
	}
}

var mutex *sync.Mutex

//登录
func (manager *ClientManager) register(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	if client.UserId != "" {
		manager.Clients[client.UserId] = client
	}
}

//下线
func (manager *ClientManager) logout(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	if client.UserId != "" {
		delete(manager.Clients,client.UserId)
	}
}

func (manager *ClientManager) StartListen()  {
	for {
		select {
		case c,ok := <- manager.Login:
			if ok {
				manager.register(c)
			}
		case c,ok := <- manager.Logout:
			if ok {
				manager.register(c)
			}
		}
	}
}

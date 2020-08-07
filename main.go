package main

import (
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"jchat/logs"
	"jchat/socketservice"
	"net/http"
	"time"
)

func main() {
	initConfig()

	http.HandleFunc("/serveWs", serveWs)

	go socketservice.StartListen()

	err := http.ListenAndServe(":7978", nil)
	if err != nil {
		logs.Logger().Error(err.Error())
		return
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		logs.Logger().Error(err.Error())
		return
	}
	logs.Logger().Info("connected WebSocket", zap.Duration("time", time.Second))
	client := socketservice.NewClient(time.Now().Unix(), conn)
	go client.Read()
	go client.Write()
}

func initConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./jchat")
}

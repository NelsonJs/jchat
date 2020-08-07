package models

import "encoding/json"

type Response struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func GetResponseByte(code int64, msg string) ([]byte, error) {
	r := Response{Code: code, Msg: msg}
	return json.Marshal(&r)
}

package models

const (
	StatusMsgTypeError = 10000
)

var statusTxt = map[int64]string{
	StatusMsgTypeError: "Message Unmarshal Error",
}

func StatusCodeTxt(code int64) string {
	return statusTxt[code]
}


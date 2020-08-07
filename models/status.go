package models

const (
	StatusMsgTypeError             = 10000
	StatusParametersError          = 10001
	StatusLoginError               = 10002
	StatusLogoutError              = 10003
	StatusNotLoginError            = 10004
	StatusReceiveUserNotLoginError = 10005

	StatusLoginSuccess  = 20000
	StatusLogoutSuccess = 20001
)

var statusTxt = map[int64]string{
	StatusMsgTypeError:             "Message Unmarshal Error",
	StatusParametersError:          "Parameters Error",
	StatusLoginError:               "login error",
	StatusLogoutError:              "logout error",
	StatusNotLoginError:            "User Not Login",
	StatusReceiveUserNotLoginError: "Receive User Not Login",

	StatusLoginSuccess:  "login success",
	StatusLogoutSuccess: "logout success",
}

func StatusCodeTxt(code int64) string {
	return statusTxt[code]
}

package errors

import (
	"net/http"
)

type Errno struct {
	HttpCode int
	Code     int
	Message  string
}

var ErrSaveTomlFileNot = &Errno{HttpCode: http.StatusInternalServerError, Code: 10001, Message: "配置文件保存失败"}

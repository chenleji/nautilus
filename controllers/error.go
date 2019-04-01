package controllers

import "net/http"

var (
	NoErr        = CtlResp{http.StatusOK, "OK", nil}
	Err404       = CtlResp{http.StatusNotFound, "资源不存在", nil}
	ErrInputData = CtlResp{http.StatusBadRequest, "输入参数异常", nil}
)

type CtlResp struct {
	HttpCode int
	Message  string
	Body     interface{}
}

func (r CtlResp) SetCode(code int) (CtlResp) {
	r.HttpCode = code
	return r
}

func (r CtlResp) SetMsg(msg string) (CtlResp) {
	r.Message = msg
	return r
}

func (r CtlResp) SetData(data interface{}) (CtlResp) {
	r.Body = data
	return r
}

type ServerResp struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

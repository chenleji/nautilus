package controllers

import (
	"encoding/json"
	"errors"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"net/http"
)

type BaseController struct {
	beego.Controller
}

func (base *BaseController) VerifyInputData(obj interface{}) (err error) {
	valid := validation.Validation{}
	ok, err := valid.Valid(obj)
	if err != nil {
		return err
	}
	if !ok {
		str := ""
		for _, err := range valid.Errors {
			str += err.Key + ":" + err.Message + ";"
		}
		return errors.New(str)
	}

	return nil
}

func (base *BaseController) RespJson(ctlErr CtlResp) {
	base.Ctx.Output.Header("content-type", "application/json")

	// default value
	if ctlErr.HttpCode == 0 {
		ctlErr.HttpCode = 200
	}

	if ctlErr.Message == "" {
		ctlErr.Message = "OK"
	}

	// construct
	code := ""
	if ctlErr.HttpCode >= http.StatusOK && ctlErr.HttpCode < http.StatusMultipleChoices {
		code = fmt.Sprintf("S%d", ctlErr.HttpCode)
	} else {
		code = fmt.Sprintf("E%d", ctlErr.HttpCode)
	}

	resp := &ServerResp{
		Code:    code,
		Message: ctlErr.Message,
	}

	s, ok := ctlErr.Body.(string)
	if ok {
		resp.Data = json.RawMessage(s)
	} else {
		resp.Data = ctlErr.Body
	}

	errBytes, _ := json.Marshal(resp)
	base.CustomAbort(http.StatusOK, string(errBytes))
}

func (base *BaseController) RespInternalError(err error) {
	resp := CtlResp{}.SetCode(http.StatusInternalServerError).SetMsg(err.Error()).SetData(err)
	base.RespJson(resp)
}

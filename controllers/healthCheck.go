package controllers

import (
	"fmt"
	"github.com/chenleji/nautilus/helper"
	"net/http"
)

// Operations about object
type HealthCheckController struct {
	BaseController
}

// @Title getEpBySrv
// @Description get dolphin health status
// @Success 200 {string}  active
// @router /health [get]
func (c *HealthCheckController) Get() {
	tool := helper.Utils{}
	if err := tool.SystemHealth(); err != nil {
		data := map[string]interface{}{
			"active":  false,
			"details": fmt.Sprintf("%s", err.Error()),
		}
		c.RespJson(CtlResp{}.
			SetCode(http.StatusInternalServerError).
			SetMsg("exception").
			SetData(data))

		return
	}

	data := map[string]interface{}{
		"active":  true,
		"details": "",
	}
	c.RespJson(CtlResp{}.SetData(data))

	return
}

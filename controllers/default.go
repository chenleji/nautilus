package controllers

import (
	"github.com/astaxie/beego"
)

// Operations about home page
// Operations about home page
type MainController struct {
	beego.Controller
}

// @Title getEpBySrv
// @Description get nautilus home page
// @Success 200 {string}  info
// @router / [get]
func (c *MainController) Get() {
	c.Data["Website"] = "xxx.xxx.com"
	c.Data["Email"] = "xxx@gmail.com"
	c.TplName = "index.tpl"
}
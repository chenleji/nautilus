package routers

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

const (
	TenantIdHeader = "yh-tenant-id"
)

var FilterUserAgent = func(ctx *context.Context) {
	if ctx.Input.Header(TenantIdHeader) == "" {
		//	ctx.Abort(http.StatusForbidden, "Invalid User-Agent!")
		logs.Info("filter working ...")
	}
}

// @APIVersion 1.0.0
// @Title nautilus Rest API
// @Description nautilus Rest API
// @Contact admin@gmail.com
// @TermsOfServiceUrl http://x.x.x.x:8010/
package routers

import (
	"github.com/astaxie/beego"
	"github.com/chenleji/nautilus/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/health", &controllers.HealthCheckController{})

	//// add your controllers
	//ns := beego.NewNamespace("/v1/orchestration",
	//	beego.NSNamespace("/project",
	//		beego.NSInclude(&controllers.ProjectOrchestrationController{}),
	//	),
	//	beego.NSNamespace("/project/service",
	//		beego.NSInclude(&controllers.ProjectServiceController{}),
	//	),
	//	beego.NSNamespace("/service",
	//		beego.NSInclude(&controllers.ServiceController{}),
	//	),
	//	beego.NSNamespace("/stack",
	//		beego.NSInclude(&controllers.StackController{}),
	//	),
	//)
	//
	//beego.AddNamespace(ns)

	// Insert filter for Auth
	beego.InsertFilter("/v1/*", beego.BeforeRouter, FilterUserAgent)
}

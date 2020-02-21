package routers

import (
	"ds-yibasuo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 登录
	beego.Router("/api/v1/login", &controllers.LoginController{}, "post:Login")
	//beego.Router("/api/v1/logout", &controllers.LoginController{}, "get:Logout")

	// 首页

	// 集群管理
	beego.Router("/api/v1/executeCluster", &controllers.ClusterController{}, "POST:ExecuteCluster")
	beego.Router("/api/v1/readLog", &controllers.ClusterController{}, "get:ReadLog")
	beego.Router("/api/v1/executeResultSignal", &controllers.ClusterController{}, "get:ExecuteResultSignal")

	// 配置管理

	// 主机管理
	beego.Router("/api/v1/createHost", &controllers.HostController{}, "post:CreateHost")
	beego.Router("/api/v1/queryHostList", &controllers.HostController{}, "get:QueryHostList")

	// 系统设置
}

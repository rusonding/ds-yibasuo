package routers

import (
	"ds-yibasuo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//主页
	//beego.Router("/", &controllers.IndexController{})
	//beego.Router("/*", &controllers.IndexController{}) //支持vue的路由

	// 登录
	beego.Router("/api/v1/login", &controllers.LoginController{}, "post:Login")
	//beego.Router("/api/v1/logout", &controllers.LoginController{}, "get:Logout")

	// 首页

	// 集群管理
	beego.Router("/api/v1/createUpdateCluster", &controllers.ClusterController{}, "POST:CreateUpdateCluster")
	beego.Router("/api/v1/deleteCluster", &controllers.ClusterController{}, "POST:DeleteCluster")
	beego.Router("/api/v1/selectCluster", &controllers.ClusterController{}, "POST:SelectCluster")
	beego.Router("/api/v1/selectClusterList", &controllers.ClusterController{}, "GET:SelectClusterList")
	beego.Router("/api/v1/executeCluster", &controllers.ClusterController{}, "POST:ExecuteCluster")
	beego.Router("/api/v1/readLog", &controllers.ClusterController{}, "get:ReadLog")
	beego.Router("/api/v1/executeResultSignal", &controllers.ClusterController{}, "get:ExecuteResultSignal")

	// 配置管理

	// 主机管理
	beego.Router("/api/v1/createHost", &controllers.HostController{}, "post:CreateHost")
	beego.Router("/api/v1/deleteHost", &controllers.HostController{}, "post:DeleteHost")
	beego.Router("/api/v1/updateHost", &controllers.HostController{}, "post:UpdateHost")
	beego.Router("/api/v1/selectHost", &controllers.HostController{}, "post:SelectHost")
	beego.Router("/api/v1/selectHostList", &controllers.HostController{}, "get:SelectHostList")

	// 系统设置
}

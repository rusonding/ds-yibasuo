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

  // 配置管理

  // 主机管理
  beego.Router("/api/v1/createHost", &controllers.HostController{}, "post:CreateHost")
  beego.Router("/api/v1/queryHostList", &controllers.HostController{}, "get:QueryHostList")

  // 系统设置

  // 部署
  beego.Router("/api/v1/deploy", &controllers.DevopsController{}, "get:Deploy")
  beego.Router("/api/v1/readLog", &controllers.DevopsController{}, "get:ReadLog")
}
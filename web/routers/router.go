package routers

import "github.com/astaxie/beego"

func init() {
  // 登录
  beego.Router("/api/v1/login", &controllers.LoginController{}, "post:Login")
  beego.Router("/api/v1/logout", &controllers.LoginController{}, "get:Logout")

  // 首页

  // 集群管理

  // 配置管理

  // 设备管理

  // 设置
}
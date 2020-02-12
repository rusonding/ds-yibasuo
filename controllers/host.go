package controllers

import (
  "ds-yibasuo/models"
  "encoding/json"
  "github.com/astaxie/beego"
  "github.com/astaxie/beego/logs"
)

type HostController struct {
  beego.Controller
}

// 创建主机
func (c *HostController) CreateHost() {
  var req models.HostInfo

  if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
    if _, err := req.HostInsert(); err == nil {
      c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
    } else {
      logs.Info(err)
    }
  } else {
    c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
  }

  c.ServeJSON()
}

// 删除主机

// 修改主机

// 查询主机列表
func (c *HostController) QueryHostList() {
  res, err := models.QueryHostList()
  if err == nil {
    c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
  } else {
    c.Data["json"] = models.Response{Code: 500, Message: "查询错误", Result: nil}
  }

  c.ServeJSON()
}

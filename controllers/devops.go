package controllers

import (
  "ds-yibasuo/models"
  "ds-yibasuo/utils"
  "github.com/astaxie/beego"
  "github.com/astaxie/beego/logs"
  "strconv"
)

type DevopsController struct {
  beego.Controller
}

// 部署
func (c *DevopsController) Deploy() {
  now := utils.Now()
  ansible := models.DevopsInfo{ExecTime: now}

  // 备份日志
  ansible.BackupLog()

  // 部署
  ansible.Deploy()

  c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
  c.ServeJSON()
}

// 读ansible执行日志
func (c *DevopsController) ReadLog() {
  now := utils.Now()
  ansible := models.DevopsInfo{ExecTime: now}

  // 分页读日志，一页10条
  page, _ := strconv.Atoi(c.Input().Get("page"))
  log, err := ansible.ReadLog((page*10)-9, page*10)
  if err != nil {
    logs.Error(err)
  }

  // 获得日志总行数
  rows, err := ansible.GetLogRows()
  if err != nil {
    logs.Error(err)
  }

  // 组装响应
  res := models.DevopsLogResult{
    CurrentPage: page,
    Rows:        rows,
    Data:        log,
  }

  c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
  c.ServeJSON()
}

package controllers

import (
	"ds-yibasuo/models"
	. "ds-yibasuo/utils/common"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type HostController struct {
	beego.Controller
}

// 创建主机
func (c *HostController) CreateHost() {
	now := Now()
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		if _, err := req.CreateHost(); err == nil {
			// ansible刷新主机
			ansible := models.DevopsInfo{ExecTime: now}
			ansible.BackupLog()
			err = ansible.RefreshHost(req.Root)
			if err != nil {
				logs.Error("refresh host err: ", err)
				c.Data["json"] = models.Response{Code: 500, Message: err.Error(), Result: nil}
			} else {
				c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
			}
		} else {
			logs.Info(err)
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// 删除主机
func (c *HostController) UpdateHost() {

}

// 修改主机
func (c *HostController) DeleteHost() {

}

// 查询具体主机
func (c *HostController) SelectHost() {

}

// 查询主机列表
func (c *HostController) SelectHostList() {
	res, err := models.SelectHostList()
	if err == nil {
		c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "查询错误", Result: nil}
	}

	c.ServeJSON()
}

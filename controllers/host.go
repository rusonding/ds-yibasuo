package controllers

import (
	"ds-yibasuo/models"
	"ds-yibasuo/utils/common"
	"ds-yibasuo/utils/ini"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type HostController struct {
	beego.Controller
}

// controller层
// 创建主机
// 验证结构体并调用model的实际方法
// 创建后并调用刷新ansible host
func (c *HostController) CreateHost() {
	logs.Info("controller create host")
	now := common.Now()
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		if err := req.CreateHost(); err == nil {
			// ansible读host.ini配置
			hosts := ini.IniHosts{}
			err := hosts.ReadHosts()
			if err != nil {
				logs.Error(err)
				return
			}
			// ansible刷新host.ini配置
			// TODO 可能有不成功的隐患
			for _, value := range hosts.Servers {
				if value != req.Ip {
					hosts.Servers = append(hosts.Servers, req.Ip)
				}
			}
			hosts.AnsibleSshPass = req.Root
			err = hosts.WriteHosts()
			if err != nil {
				logs.Error(err)
				return
			}
			// ansible刷新主机
			dev := models.DevopsInfo{ExecTime: now}
			dev.BackupLog(models.DeployUpdate)
			err = dev.RefreshHost(req.Root)
			if err != nil {
				logs.Error("refresh host err: ", err)
				c.Data["json"] = models.Response{Code: 500, Message: err.Error(), Result: nil}
			} else {
				c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
			}
		} else {
			logs.Info(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("创建错误: %s", err), Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 删除主机
func (c *HostController) DeleteHost() {
	logs.Info("controller delete host")
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		err := req.DeleteHost()
		if err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("删除发生错误！%s", err), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 修改主机
func (c *HostController) UpdateHost() {
	logs.Info("controller update host")
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		err := req.UpdateHost()
		if err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("修改发生错误！%s", err), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 查询具体主机
func (c *HostController) SelectHost() {
	logs.Info("controller select host")
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		res, err := req.SelectHost()
		if err != nil {
			logs.Error(err)
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

// controller层
// 查询主机列表
func (c *HostController) SelectHostList() {
	logs.Info("controller select host list")
	page, err := strconv.Atoi(c.Input().Get("page"))
	if err != nil {
		c.Data["json"] = models.Response{Code: 500, Message: "查询错误", Result: nil}
	} else {
		res, err := models.SelectHostList(page)
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		}
	}

	c.ServeJSON()
}

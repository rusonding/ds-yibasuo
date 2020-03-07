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
// TODO 如果已经成功建立用户了,再去修改root密码,会不成功
func (c *HostController) CreateHost() {
	logs.Info("controller create host")
	now := common.Now()
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 判断重名
		if ok, err := req.CheckName(); err == nil {
			if ok {
				c.Data["json"] = models.Response{Code: 500, Message: "名称冲突", Result: nil}
				c.ServeJSON()
				return
			}
		}

		// ansible刷新host.ini配置
		hosts := ini.IniHosts{}
		hostFuck := make(map[string]string)
		hostFuck["ip"] = req.Ip
		hostFuck["pwd"] = req.Root
		hostFuck["port"] = strconv.Itoa(req.Port)
		hosts.Servers = append(hosts.Servers, hostFuck)
		hosts.AnsibleSshPass = req.Root
		hosts.AnsibleUser = "easy" //TODO 暂时写死 用户只能叫easy
		if err = hosts.WriteHosts(); err != nil {
			logs.Error("host ini write err: ", err)
		}
		// ansible刷新主机
		dev := models.DevopsInfo{ExecTime: now}
		dev.BackupLog(models.DeployUpdate)
		if err = dev.RefreshHost(req.Root); err != nil {
			logs.Error("refresh host err: ", err)
			c.Data["json"] = models.Response{Code: 500, Message: err.Error(), Result: nil}
			c.ServeJSON()
			return
		}
		// 刷新成功则开始创建host
		if err := req.CreateHost(); err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
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
// TODO 如果已经成功建立用户了,再去修改root密码,会不成功
func (c *HostController) UpdateHost() {
	logs.Info("controller update host")
	var req models.HostInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 判断重名
		if ok, err := req.CheckName(); err == nil {
			if ok {
				c.Data["json"] = models.Response{Code: 500, Message: "名称冲突", Result: nil}
				c.ServeJSON()
				return
			}
		}

		if req.Id == "" {
			c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
			c.ServeJSON()
			return
		}
		// ansible刷新host.ini配置
		hosts := ini.IniHosts{}
		hostFuck := make(map[string]string)
		hostFuck["ip"] = req.Ip
		hostFuck["pwd"] = req.Root
		hostFuck["port"] = strconv.Itoa(req.Port)
		hosts.Servers = append(hosts.Servers, hostFuck)
		hosts.AnsibleSshPass = req.Root
		hosts.AnsibleUser = "easy" //TODO 暂时写死 用户只能叫easy
		if err = hosts.WriteHosts(); err != nil {
			logs.Error("host ini write err: ", err)
		}
		// ansible刷新主机
		dev := models.DevopsInfo{}
		dev.BackupLog(models.DeployUpdate)
		if err = dev.RefreshHost(req.Root); err != nil {
			logs.Error("refresh host err: ", err)
			c.Data["json"] = models.Response{Code: 500, Message: err.Error(), Result: nil}
			c.ServeJSON()
			return
		}
		// 刷新成功再修改
		if err := req.UpdateHost(); err != nil {
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
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
		} else if err.Error() == "null" {
			c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
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
		} else if err.Error() == "null" {
			c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		}
	}

	c.ServeJSON()
}

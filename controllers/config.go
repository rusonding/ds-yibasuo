package controllers

import (
	"ds-yibasuo/models"
	"ds-yibasuo/utils/black"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type ConfigController struct {
	beego.Controller
}

func (c *ConfigController) CreateConfig() {
	logs.Info("controller create config")
	var req models.ConfigInfo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 类型判断
		if ok, err := black.Contain(models.ConfigList, req.Typ); err != nil && !ok {
			c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
			c.ServeJSON()
			return
		}
		if err := req.CreateConfig(); err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("插入错误：%s", err), Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

func (c *ConfigController) DeleteConfig() {
	logs.Info("controller delete config")
	var req models.ConfigInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		err := req.DeleteConfig()
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

func (c *ConfigController) UpdateConfig() {
	logs.Info("controller update config")
	var req models.ConfigInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		if req.Id == "" {
			c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
			c.ServeJSON()
			return
		}
		if err := req.UpdateConfig(); err != nil {
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

func (c *ConfigController) SelectConfig() {
	logs.Info("controller select config")
	var req models.ConfigInfo

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		res, err := req.SelectConfig()
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

func (c *ConfigController) SelectConfigList() {
	logs.Info("controller select config list")

	page, err := strconv.Atoi(c.Input().Get("page"))
	typ := c.Input().Get("typ")
	// 类型判断
	if ok, err := black.Contain(models.ConfigList, typ); err != nil && !ok {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
		c.ServeJSON()
		return
	}
	if err != nil {
		c.Data["json"] = models.Response{Code: 500, Message: "查询错误", Result: nil}
	} else {
		res, err := models.SelectConfigList(page, typ)
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: res}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		}
	}

	c.ServeJSON()
}

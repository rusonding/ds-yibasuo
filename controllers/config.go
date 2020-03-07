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
	var req models.ConfigInfoPuppet
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		// 类型判断
		if ok, err := black.Contain(models.ConfigList, req.Typ); err != nil && !ok {
			c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
			c.ServeJSON()
			return
		}
		var conf models.ConfigInfo
		switch req.Typ {
		case "frontend":
			conf.Conf = &models.ConfigFrontend{Frontend: &models.Frontends{EscProxyPort: req.EscProxyPort}}
		case "backend":
			conf.Conf = &models.ConfigBackend{Backend: &models.Backends{ServerPort: req.ServerPort}}
		case "alert":
			conf.Conf = &models.ConfigAlert{Alert: models.AlertMail{
				AlertType:    req.AlertType,
				MailProtocol: req.MailProtocol,
				MailSender:   req.MailSender,
				MailUser:     req.MailUser,
				MailPasswd:   req.MailPasswd,
				MailSmtpHost: req.MailSmtpHost,
				MailSmtpPort: req.MailSmtpPort,
				MailSsl:      req.MailSsl,
				MailTls:      req.MailTls,
			}}
		case "master":
			conf.Conf = &models.ConfigMaster{Master: &models.Masters{MasterReservedMemory: req.MasterReservedMemory}}
		case "worker":
			conf.Conf = &models.ConfigWorker{Worker: &models.Workers{WorkerReservedMemory: req.WorkerReservedMemory}}
		case "database":
			conf.Conf = &models.ConfigDatabase{
				DatabaseType: req.DatabaseType,
				DatabaseName: req.DatabaseName,
				Account:      req.Account,
				Password:     req.Password,
			}
		case "hadoop":
			conf.Conf = &models.ConfigHadoop{Hadoop: &models.Hadoops{FsDefaultFS: req.FsDefaultFS}}
		case "common":
			conf.Conf = &models.ConfigCommon{Common: &models.Commons{DataStore2hdfsBasepath: req.DataStore2hdfsBasepath}}
		}
		conf.Name = req.Name
		conf.Typ = req.Typ
		conf.Remark = req.Remark
		//判断重名
		if ok, err := conf.CheckName(); err == nil {
			if ok {
				c.Data["json"] = models.Response{Code: 500, Message: "名称冲突", Result: nil}
				c.ServeJSON()
				return
			}
		}
		//开始插入
		if err := conf.CreateConfig(); err == nil {
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
	var req models.ConfigInfoPuppet

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err == nil {
		var conf models.ConfigInfo
		switch req.Typ {
		case "frontend":
			conf.Conf = &models.ConfigFrontend{Frontend: &models.Frontends{EscProxyPort: req.EscProxyPort}}
		case "backend":
			conf.Conf = &models.ConfigBackend{Backend: &models.Backends{ServerPort: req.ServerPort}}
		case "alert":
			conf.Conf = &models.ConfigAlert{Alert: models.AlertMail{
				AlertType:    req.AlertType,
				MailProtocol: req.MailProtocol,
				MailSender:   req.MailSender,
				MailUser:     req.MailUser,
				MailPasswd:   req.MailPasswd,
				MailSmtpHost: req.MailSmtpHost,
				MailSmtpPort: req.MailSmtpPort,
				MailSsl:      req.MailSsl,
				MailTls:      req.MailTls,
			}}
		case "master":
			conf.Conf = &models.ConfigMaster{Master: &models.Masters{MasterReservedMemory: req.MasterReservedMemory}}
		case "worker":
			conf.Conf = &models.ConfigWorker{Worker: &models.Workers{WorkerReservedMemory: req.WorkerReservedMemory}}
		case "database":
			conf.Conf = &models.ConfigDatabase{
				DatabaseType: req.DatabaseType,
				DatabaseName: req.DatabaseName,
				Account:      req.Account,
				Password:     req.Password,
			}
		case "hadoop":
			conf.Conf = &models.ConfigHadoop{Hadoop: &models.Hadoops{FsDefaultFS: req.FsDefaultFS}}
		case "common":
			conf.Conf = &models.ConfigCommon{Common: &models.Commons{DataStore2hdfsBasepath: req.DataStore2hdfsBasepath}}
		}
		conf.Id = req.Id
		conf.Name = req.Name
		conf.Typ = req.Typ
		conf.Remark = req.Remark

		if req.Id == "" {
			c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
			c.ServeJSON()
			return
		}
		//判断重名
		if ok, err := conf.CheckName(); err == nil {
			if ok {
				c.Data["json"] = models.Response{Code: 500, Message: "名称冲突", Result: nil}
				c.ServeJSON()
				return
			}
		}
		if err := conf.UpdateConfig(); err != nil {
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

func (c *ConfigController) SelectConfigList() {
	logs.Info("controller select config list")

	page, err := strconv.Atoi(c.Input().Get("page"))
	typ := c.Input().Get("typ")
	// 获取全部配置
	if typ == "all" {
		all, err := models.SelectAllConfig()
		if err == nil {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: all}
			c.ServeJSON()
			return
		} else if err.Error() == "null" {
			c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
			c.ServeJSON()
			return
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
			c.ServeJSON()
			return
		}
	}

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
		} else if err.Error() == "null" {
			c.Data["json"] = models.Response{Code: 200, Message: "没有查到", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("查询发生错误！%s", err), Result: nil}
		}
	}

	c.ServeJSON()
}

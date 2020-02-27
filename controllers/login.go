package controllers

import (
	"ds-yibasuo/models"
	"ds-yibasuo/utils/black"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Login() {
	logs.Info("controller login: \n", black.Byte2String(c.Ctx.Input.RequestBody))
	var user models.User

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err == nil {
		if ok, err := user.UserCheck(user.Password); ok {
			c.Data["json"] = models.Response{Code: 200, Message: "ok", Result: nil}
		} else {
			c.Data["json"] = models.Response{Code: 500, Message: fmt.Sprintf("登陆错误：%s", err), Result: nil}
		}
	} else {
		c.Data["json"] = models.Response{Code: 500, Message: "参数错误", Result: nil}
	}

	c.ServeJSON()
}

package controllers

import (
  "github.com/astaxie/beego"
  . "web/models"
)

type LoginController struct {
  beego.Controller
}

func (c *LoginController) Login()  {
  var user User

  if response, err := user.UserCheck(c.Ctx.Input.RequestBody); err != nil {
    c.Data["json"] = response
    return
  }


}
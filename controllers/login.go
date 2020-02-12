package controllers

import (
  . "ds-yibasuo/models"
  "encoding/json"
  "github.com/astaxie/beego"
)

type LoginController struct {
  beego.Controller
}

func (c *LoginController) Login() {
  var user User

  if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
    c.Data["json"] = Response{Code: 500, Message: "参数错误", Result: nil}
    c.ServeJSON()
    return
  } else {
    c.Data["json"] = Response{Code: 200, Message: user.Password, Result: nil}
    c.ServeJSON()
    return
  }
}

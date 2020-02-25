package controllers

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	c.Data["Website"] = "ds-yibasuo"
	c.TplName = "index.html"
	_ = c.Render()
}

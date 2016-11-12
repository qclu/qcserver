package controllers

import (
	"github.com/astaxie/beego"
)

type QcErrorController struct {
	beego.Controller
}

func (c *QcErrorController) Error501() {
	c.Data["content"] = "server error"
	c.TplName = "501.tpl"
}

func (c *QcErrorController) ErrorDbFailed() {
	c.Data["content"] = "database is now down"
	c.TplName = "dberror.tpl"
}

func (c *QcErrorController) ErrorReqInvalid() {
	c.Data["content"] = "database is now down"
	c.TplName = "dberror.tpl"
}

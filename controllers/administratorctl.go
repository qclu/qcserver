package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcAdministratorCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcAdministratorCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcAdministratorCtl) Post() {
	var ob models.QcAdministrator
	var err error
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcAdmin(o.dbSync, ob.Username, ob.Password, ob.Role)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcAdministratorCtl) Delete() {
	hname := h.GetString("name")
	err := h.dbSync.DeleteQcAdminSQL(hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

// @router / [get]
func (h *QcAdministratorCtl) Get() {
	hname := h.GetString("name")
	admin, err := h.dbSync.GetQcAdmin(hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = admin
	h.ServeJSON()
}

// @router /list [get]
func (h *QcAdministratorCtl) GetList() {
	pgidx_str := h.GetString("pageidx")
	h.logger.LogInfo("pageidx: ", pgidx_str)
	pgidx, err := strconv.Atoi(pgidx_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pageidx' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pageidx' from request, err: " + err.Error()
		h.ServeJSON()
	}
	pgsize_str := h.GetString("pagesize")
	h.logger.LogInfo("pagesize: ", pgsize_str)
	pgsize, err := strconv.Atoi(pgsize_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pagesize' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
	}
	h.logger.LogInfo("list hospital info(pageidx: ", pgidx, ", pagesize: ", pgsize, ")")
	admins, err := h.dbSync.GetQcAdmins(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = admins
	h.ServeJSON()
}

// @router / [PUT]
func (h *QcAdministratorCtl) Update() {
	hname := h.GetString("org_name")
	if len(hname) == 0 {
		h.logger.LogError("failed to parse admin name from request")
		h.Data["json"] = "failed to parse admin name from request"
		h.ServeJSON()
	}
	admin, err := h.dbSync.GetQcAdmin(hname)
	if err != nil {
		h.logger.LogError("failed to get admin[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get admin[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		admin.Username = new_name
	}
	new_pwd := h.GetString("passwd")
	if len(new_pwd) > 0 {
		admin.Password = new_pwd
	}
	new_role_str := h.GetString("role")
	if len(new_role_str) > 0 {
		new_role, err := strconv.Atoi(new_role_str)
		if err != nil {
			h.logger.LogError("invalid parameter value for role, err: ", err)
			h.Data["json"] = "invalid parameter value for role, err: " + err.Error()
			h.ServeJSON()
		}
		admin.Role = new_role
	}
	err = h.dbSync.UpdateQcAdmin(admin)
	if err != nil {
		h.logger.LogError("failed to update hospital[", hname, "], err: ", err)
		h.Data["json"] = "failed to update hospital[" + hname + "], err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = admin
	h.ServeJSON()
}

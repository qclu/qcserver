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
	o.logger.LogInfo("request body:", string(o.Ctx.Input.RequestBody))
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	o.logger.LogInfo("username: ", ob.Username)
	pob, err := models.CreateQcAdmin(o.dbSync, ob.Username, ob.Password, ob.Role)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
		return
	}
	o.Data["json"] = pob
	o.ServeJSON()
	return
}

// @router / [delete]
func (h *QcAdministratorCtl) Delete() {
	idstr := h.GetString("id")
	err := h.dbSync.DeleteQcObjectWithID(idstr, models.DB_T_ADMINISTRATOR)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
	return
}

// @router / [get]
func (h *QcAdministratorCtl) Get() {
	hname := h.GetString("name")
	idstr := h.GetString("id")
	var admin *models.QcAdministrator
	var err error
	if len(hname) > 0 {
		admin, err = h.dbSync.GetQcAdmin(hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		h.logger.LogInfo("id: ", idstr)
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get admin")
			h.Data["json"] = "invalid id value to get admin"
			h.ServeJSON()
			return
		}
		admin, err = h.dbSync.GetQcAdminWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for admin get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = admin
	h.ServeJSON()
	return
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
		return
	}
	pgsize_str := h.GetString("pagesize")
	h.logger.LogInfo("pagesize: ", pgsize_str)
	pgsize, err := strconv.Atoi(pgsize_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pagesize' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.logger.LogInfo("list hospital info(pageidx: ", pgidx, ", pagesize: ", pgsize, ")")
	admins, err := h.dbSync.GetQcAdmins(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_ADMINISTRATOR)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(admins)
	res_data["objects"] = admins
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcAdministratorCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse admin id from request")
		h.Data["json"] = "failed to parse admin id from request"
		h.ServeJSON()
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		h.logger.LogError("invalid value for admin id from request")
		h.Data["json"] = "invalid value for admin id from request"
		h.ServeJSON()
		return
	}
	admin, err := h.dbSync.GetQcAdminWithId(id)
	if err != nil {
		h.logger.LogError("failed to get admin[", id, "] from database, err: ", err)
		h.Data["json"] = "failed to get admin[" + idstr + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
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
			return
		}
		admin.Role = new_role
	}
	err = h.dbSync.UpdateQcAdmin(admin)
	if err != nil {
		h.logger.LogError("failed to update hospital[", idstr, "], err: ", err)
		h.Data["json"] = "failed to update hospital[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = admin
	h.ServeJSON()
	return
}

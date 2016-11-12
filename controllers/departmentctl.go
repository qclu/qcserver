package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	//"strconv"
)

// Operations about object
type QcDepartmentCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcDepartmentCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @Title Create
// @router / [post]
func (o *QcDepartmentCtl) Post() {
	hname := o.GetString("hname")
	hospital, err := o.dbSync.GetQcHospital(hname)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get hospital[" + hname + "] info, err: " + err.Error()
		o.ServeJSON()
	}
	var ob models.QcDepartment
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	o.logger.LogInfo("Department create request: ", ob)
	pob, err := models.CreateQcDepartment(o.dbSync, ob.Name, hospital)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcDepartmentCtl) Delete() {
	hname := h.GetString("hname")
	dname := h.GetString("dname")
	err := h.dbSync.DeleteQcDepartmentSQL(hname, dname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

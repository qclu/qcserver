package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
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
	h.logger.LogInfo("delete department[", dname, "] from hospital[", hname, "]")
	err := models.DeleteQcDepartment(h.dbSync, dname, hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = h
	h.ServeJSON()
}

// @router / [get]
func (h *QcDepartmentCtl) Get() {
	hname := h.GetString("name")
	h.logger.LogInfo("get hospital info, name: ", hname)
	hospital, err := h.dbSync.GetQcHospital(hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = hospital
	h.ServeJSON()
}

// @router /list [get]
func (h *QcDepartmentCtl) GetList() {
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
	hospitals, err := h.dbSync.GetQcHospitals(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = hospitals
	h.ServeJSON()
}

// @router / [PUT]
func (h *QcDepartmentCtl) Update() {
	hname := h.GetString("org_name")
	h.logger.LogInfo("update hospital ", hname)
	if len(hname) == 0 {
		h.logger.LogError("failed to parse hospital name from request")
		h.Data["json"] = "failed to parse hospital name from request"
		h.ServeJSON()
	}
	hospital, err := h.dbSync.GetQcHospital(hname)
	if err != nil {
		h.logger.LogError("failed to get hospital[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get hospital[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		hospital.Name = new_name
	}
	new_gis := h.GetString("gis")
	if len(new_gis) > 0 {
		hospital.Gis = new_gis
	}
	new_addr := h.GetString("addr")
	if len(new_addr) > 0 {
		hospital.Addr = new_addr
	}
	err = h.dbSync.UpdateQcHospital(hospital)
	if err != nil {
		h.logger.LogError("failed to update hospital[", hname, "], err: ", err)
		h.Data["json"] = "failed to update hospital[" + hname + "], err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = hospital
	h.ServeJSON()
}

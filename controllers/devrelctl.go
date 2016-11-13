package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcDevRelCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcDevRelCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcDevRelCtl) Post() {
	swversion_str := o.GetString("swversion")
	swv, err := o.dbSync.GetQcSwVersion(swversion_str)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get swversion[" + swversion_str + "] info, err: " + err.Error()
		o.ServeJSON()
	}
	hospital_str := o.GetString("hospital")
	department_str := o.GetString("department")
	department, err := o.dbSync.GetQcDepartment(department_str, hospital_str)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get department[" + department_str + ":" + hospital_str + "] info, err: " + err.Error()
		o.ServeJSON()
	}
	o.logger.LogInfo("department info: ", department)
	var devrel models.QcDevRel
	json.Unmarshal(o.Ctx.Input.RequestBody, &devrel)
	pob, err := models.CreateQcDevRel(o.dbSync, devrel.Sn, devrel.SmCard, devrel.Date, swv, department)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcDevRelCtl) Delete() {
	sn := h.GetString("sn")
	err := h.dbSync.DeleteQcDevRelSQL(sn)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

// @router / [get]
func (h *QcDevRelCtl) Get() {
	sn := h.GetString("sn")
	devrel, err := h.dbSync.GetQcDevRel(sn)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = devrel
	h.ServeJSON()
}

// @router /list [get]
func (h *QcDevRelCtl) GetList() {
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
	devrels, err := h.dbSync.GetQcDevRels(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get dev rels list, err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = devrels
	h.ServeJSON()
}

// @router / [PUT]
func (h *QcDevRelCtl) Update() {
	sn := h.GetString("org_sn")
	if len(sn) == 0 {
		h.logger.LogError("failed to parse sn info from request")
		h.Data["json"] = "failed to parse sn info from request"
		h.ServeJSON()
	}
	devrel, err := h.dbSync.GetQcDevRel(sn)
	if err != nil {
		h.logger.LogError("failed to get devrel[sn: ", sn, "] from database, err: ", err)
		h.Data["json"] = "failed to get devrel[sn:" + sn + "] from database, err: " + err.Error()
		h.ServeJSON()
	}
	new_sn := h.GetString("sn")
	if len(new_sn) > 0 {
		devrel.Sn = new_sn
	}
	new_date := h.GetString("date")
	if len(new_date) > 0 {
		devrel.Date = new_date
	}
	new_smcard := h.GetString("smcard")
	if len(new_smcard) > 0 {
		devrel.SmCard = new_smcard
	}
	new_swv := h.GetString("swversion")
	if len(new_swv) > 0 {
		new_swv_obj, err := h.dbSync.GetQcSwVersion(new_swv)
		if err != nil {
			h.logger.LogError("failed to get swversion[", new_swv, "] info, err: ", err)
			h.Data["json"] = "failed to get swversion[" + new_swv + "] info, err: " + err.Error()
			h.ServeJSON()
		}
		devrel.SwVersion = new_swv_obj
	}
	new_hospital := h.GetString("hospital")
	new_department := h.GetString("department")
	if len(new_hospital) > 0 && len(new_department) > 0 {
		new_receiver, err := h.dbSync.GetQcDepartment(new_hospital, new_department)
		if err != nil {
			h.logger.LogError("failed to get receiver[", new_hospital, ":", new_department, "] info, err: ", err)
			h.Data["json"] = "failed to get receiver[" + new_hospital + ":" + new_department + "] info, err: " + err.Error()
			h.ServeJSON()
		}
		devrel.Receiver = new_receiver
	} else if len(new_hospital)+len(new_department) > 0 {
		h.Data["json"] = "hospital and department if one set must both are set"
		h.ServeJSON()
	}
	err = h.dbSync.UpdateQcDevRel(devrel)
	if err != nil {
		h.logger.LogError("failed to update devrel[", sn, "], err: ", err)
		h.Data["json"] = "failed to update devrel[" + sn + "], err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = devrel
	h.ServeJSON()
}
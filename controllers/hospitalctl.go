package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcHospitalCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcHospitalCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcHospitalCtl) Post() {
	var ob models.QcHospital
	var err error
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	o.logger.LogInfo("RequestBody: ", o.Ctx.Request)
	o.logger.LogInfo("Hospital create request: ", ob)
	if o.dbSync == nil {
		o.logger.LogError("dbSync is uninitialized when trying to create QcHospital, request: ", o.Ctx)
		o.Abort("501")
	}
	pob, err := models.CreateQcHospital(o.dbSync, ob.Name, ob.Addr, ob.Gis)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
		return
	}
	o.Data["json"] = map[string]models.QcHospital{"Object": *pob}
	o.ServeJSON()
	return
}

// @router / [delete]
func (h *QcHospitalCtl) Delete() {
	hname := h.GetString("name")
	h.logger.LogInfo("delete hospital info, name: ", hname)
	err := models.DeleteQcHospital(h.dbSync, hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
	return
}

// @router / [get]
func (h *QcHospitalCtl) Get() {
	hname := h.GetString("name")
	h.logger.LogInfo("get hospital info, name: ", hname)
	hospital, err := h.dbSync.GetQcHospital(hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = hospital
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcHospitalCtl) GetList() {
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
	hospitals, err := h.dbSync.GetQcHospitals(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_HOSPITAL)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = entcnt
	res_data["objects"] = hospitals
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcHospitalCtl) Update() {
	hname := h.GetString("org_name")
	h.logger.LogInfo("update hospital ", hname)
	if len(hname) == 0 {
		h.logger.LogError("failed to parse hospital name from request")
		h.Data["json"] = "failed to parse hospital name from request"
		h.ServeJSON()
		return
	}
	hospital, err := h.dbSync.GetQcHospital(hname)
	if err != nil {
		h.logger.LogError("failed to get hospital[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get hospital[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
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
		return
	}
	h.Data["json"] = hospital
	h.ServeJSON()
	return
}

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
	o.logger.LogInfo("RequestBody: ", string(o.Ctx.Input.RequestBody))
	o.logger.LogInfo("Hospital create request: ", ob)
	if o.dbSync == nil {
		o.logger.LogError("dbSync is uninitialized when trying to create QcHospital, request: ", o.Ctx)
		o.Abort("501")
	}
	pob, err := models.CreateQcHospital(o.dbSync, ob.Name, ob.Prov, ob.City, ob.Addr, ob.Gis, ob.Level)
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
	idstr := h.GetString("id")
	_, err := strconv.Atoi(idstr)
	if err != nil {
		h.logger.LogError("invalid id value to delete hospital")
		h.Data["json"] = "invalid id value to delete hospital"
		h.ServeJSON()
		return
	}
	err = h.dbSync.DeleteQcObjectWithID(idstr, models.DB_T_HOSPITAL)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "internal database operation error"
		h.ServeJSON()
		return
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
	return
}

// @router / [get]
func (h *QcHospitalCtl) Get() {
	hname := h.GetString("name")
	idstr := h.GetString("id")
	var hospital *models.QcHospital
	var err error
	if len(hname) > 0 {
		hospital, err = h.dbSync.GetQcHospital(hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Abort("501")
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get hospital")
			h.Data["json"] = "invalid id value to get hospital"
			h.ServeJSON()
			return
		}
		hospital, err = h.dbSync.GetQcHospitalWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for hospital get"
		h.ServeJSON()
		return
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
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_HOSPITAL)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(hospitals)
	res_data["objects"] = hospitals
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcHospitalCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse hospital id from request")
		h.Data["json"] = "failed to parse hospital id from request"
		h.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(idstr)
	hospital, err := h.dbSync.GetQcHospitalWithId(id)
	if err != nil {
		h.logger.LogError("failed to get hospital[", id, "] from database, err: ", err)
		h.Data["json"] = "failed to get hospital[" + idstr + "] from database, err: " + err.Error()
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
		h.logger.LogError("failed to update hospital[", id, "], err: ", err)
		h.Data["json"] = "failed to update hospital[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = hospital
	h.ServeJSON()
	return
}

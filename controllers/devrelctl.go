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
		return
	}
	hospital_str := o.GetString("hospital")
	department_str := o.GetString("department")
	department, err := o.dbSync.GetQcDepartment(department_str, hospital_str)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get department[" + department_str + ":" + hospital_str + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	o.logger.LogInfo("department info: ", department)
	var devrel models.QcDevRel
	json.Unmarshal(o.Ctx.Input.RequestBody, &devrel)
	pob, err := models.CreateQcDevRel(o.dbSync, devrel.Sn, devrel.SmCard, devrel.Date, swv, department)
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
func (h *QcDevRelCtl) Delete() {
	idstr := h.GetString("id")
	err := h.dbSync.DeleteQcObjectSQL(idstr, models.DB_T_DEVREL)
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
func (h *QcDevRelCtl) Get() {
	sn := h.GetString("sn")
	idstr := h.GetString("id")
	var devrel *models.QcDevRel
	var err error
	if len(sn) > 0 {
		devrel, err = h.dbSync.GetQcDevRel(sn)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get devrel")
			h.Data["json"] = "invalid id value to get devrel"
			h.ServeJSON()
			return
		}
		devrel, err = h.dbSync.GetQcDevRelWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for devrel get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = devrel
	h.ServeJSON()
	return
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
	devid_str := h.GetString("devid")
	if len(devid_str) > 0 {
		_, err = strconv.Atoi(devid_str)
		if err != nil {
			h.logger.LogError("failed to parse 'devid' from request, err: ", err)
			h.Data["json"] = "invalid value for 'devid' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	serial := h.GetString("serial")
	start_date := h.GetString("startdate")
	end_date := h.GetString("enddate")
	departmentid_str := h.GetString("departmentid")
	if len(departmentid_str) > 0 {
		_, err = strconv.Atoi(departmentid_str)
		if err != nil {
			h.logger.LogError("failed to parse 'departmentid' from request, err: ", err)
			h.Data["json"] = "invalid value for 'departmentid' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	hid_str := h.GetString("hospitalid")
	if len(hid_str) > 0 {
		_, err = strconv.Atoi(hid_str)
		if err != nil {
			h.logger.LogError("failed to parse 'hospitalid' from request, err: ", err)
			h.Data["json"] = "invalid value for 'hospitalid' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}

	devrels, err := h.dbSync.GetQcDevRelsCond(pgidx_str, pgsize_str, devid_str, serial, start_date, end_date, departmentid_str, hid_str)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get dev rels list, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_DEVREL)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(devrels)
	res_data["objects"] = devrels
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcDevRelCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse id info from request")
		h.Data["json"] = "failed to parse id info from request"
		h.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(idstr)
	devrel, err := h.dbSync.GetQcDevRelWithId(id)
	if err != nil {
		h.logger.LogError("failed to get devrel[id: ", idstr, "] from database, err: ", err)
		h.Data["json"] = "failed to get devrel[id:" + idstr + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
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
			return
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
			return
		}
		devrel.Receiver = new_receiver
	} else if len(new_hospital)+len(new_department) > 0 {
		h.Data["json"] = "hospital and department if one set must both are set"
		h.ServeJSON()
		return
	}
	err = h.dbSync.UpdateQcDevRel(devrel)
	if err != nil {
		h.logger.LogError("failed to update devrel[", idstr, "], err: ", err)
		h.Data["json"] = "failed to update devrel[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = devrel
	h.ServeJSON()
	return
}

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcReagentRelCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcReagentRelCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @Title Create
// @router / [post]
func (o *QcReagentRelCtl) Post() {
	dname := o.GetString("department")
	hname := o.GetString("hospital")
	department_obj, err := o.dbSync.GetQcDepartment(dname, hname)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get department[" + hname + ":" + dname + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	produceserial := o.GetString("produceserial")
	regproduce, err := o.dbSync.GetQcReagentProduce(produceserial)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get regproduce[" + produceserial + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	var ob models.QcReagentRel
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcReagentRel(o.dbSync, ob.ReleaseTime, ob.ReleaseSerial, ob.Annotation, ob.Amounts, regproduce, department_obj)
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
func (h *QcReagentRelCtl) Delete() {
	idstr := h.GetString("id")
	err := h.dbSync.DeleteQcObjectWithID(idstr, models.DB_T_REGREL)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

// @router / [get]
func (h *QcReagentRelCtl) Get() {
	serial := h.GetString("serial")
	idstr := h.GetString("id")
	var regrel *models.QcReagentRel
	var err error
	if len(serial) > 0 {
		regrel, err = h.dbSync.GetQcReagentRel(serial)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get reagent rel")
			h.Data["json"] = "invalid id value to get reagent rel"
			h.ServeJSON()
			return
		}
		regrel, err = h.dbSync.GetQcReagentRelWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for reagentrel get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = regrel
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcReagentRelCtl) GetList() {
	pgidx_str := h.GetString("pageidx")
	h.logger.LogInfo("pageidx: ", pgidx_str)
	_, err := strconv.Atoi(pgidx_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pageidx' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pageidx' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	pgsize_str := h.GetString("pagesize")
	h.logger.LogInfo("pagesize: ", pgsize_str)
	_, err = strconv.Atoi(pgsize_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pagesize' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	did_str := h.GetString("departmentid")
	if len(did_str) > 0 {
		_, err = strconv.Atoi(did_str)
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
	regmodelid_str := h.GetString("regmodelid")
	if len(regmodelid_str) > 0 {
		_, err = strconv.Atoi(regmodelid_str)
		if err != nil {
			h.logger.LogError("failed to parse 'regmodelid' from request, err: ", err)
			h.Data["json"] = "invalid value for 'regmodelid' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	start_date := h.GetString("startdate")
	end_date := h.GetString("enddate")
	regrels, err := h.dbSync.GetQcRegRelsCond(pgidx_str, pgsize_str, regmodelid_str, start_date, end_date, hid_str, did_str)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_DEPARTMENT)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(regrels)
	res_data["objects"] = regrels
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcReagentRelCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse reagentrel id from request")
		h.Data["json"] = "failed to parse reagentrel id from request"
		h.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(idstr)
	regrel, err := h.dbSync.GetQcReagentRelWithId(id)
	if err != nil {
		h.logger.LogError("failed to get regrel[", id, "] from database, err: ", err)
		h.Data["json"] = "failed to get regrel[" + idstr + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_reltime := h.GetString("reltime")
	if len(new_reltime) > 0 {
		regrel.ReleaseTime = new_reltime
	}
	new_relserial := h.GetString("serial")
	if len(new_relserial) > 0 {
		regrel.ReleaseSerial = new_relserial
	}
	new_dname := h.GetString("department")
	new_hname := h.GetString("hospital")
	if len(new_dname) > 0 && len(new_hname) > 0 {
		new_department, err := h.dbSync.GetQcDepartment(new_dname, new_hname)
		if err != nil {
			h.logger.LogError("failed to get department[", new_hname, ":", new_dname, "] info, err: ", err)
			h.Data["json"] = "failed to get hospital[" + new_hname + ":" + new_dname + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		regrel.Department = new_department
	}
	err = h.dbSync.UpdateQcReagentRel(regrel)
	if err != nil {
		h.logger.LogError("failed to update reagentrel[", id, "], err: ", err)
		h.Data["json"] = "failed to update reagentrel[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = regrel
	h.ServeJSON()
	return
}

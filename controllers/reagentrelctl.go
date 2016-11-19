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
	}
	var ob models.QcReagentRel
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcReagentRel(o.dbSync, ob.ReleaseTime, ob.ReleaseSerial, ob.CarId, ob.Annotation, department_obj)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcReagentRelCtl) Delete() {
	serial := h.GetString("serial")
	err := h.dbSync.DeleteQcReagentRelSQL(serial)
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
	regrels, err := h.dbSync.GetQcReagentRels(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_DEPARTMENT)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = entcnt
	res_data["objects"] = regrels
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcReagentRelCtl) Update() {
	serial := h.GetString("org_serial")
	if len(serial) == 0 {
		h.logger.LogError("failed to parse reagentrel serial from request")
		h.Data["json"] = "failed to parse reagentrel serial from request"
		h.ServeJSON()
		return
	}
	regrel, err := h.dbSync.GetQcReagentRel(serial)
	if err != nil {
		h.logger.LogError("failed to get regrel[", serial, "] from database, err: ", err)
		h.Data["json"] = "failed to get regrel[" + serial + "] from database, err: " + err.Error()
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
		h.logger.LogError("failed to update reagentrel[", serial, "], err: ", err)
		h.Data["json"] = "failed to update reagentrel[" + serial + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = regrel
	h.ServeJSON()
	return
}

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
		return
	}
	var ob models.QcDepartment
	err = json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	if err != nil {
		o.logger.LogError("invalid parameter, err: ", err)
		o.Data["json"] = string("invalid parameter, err:") + err.Error()
		o.ServeJSON()
		return
	}
	o.logger.LogInfo("department:", string(o.Ctx.Input.RequestBody))
	o.logger.LogInfo("Department create request: ", ob)
	pob, err := models.CreateQcDepartment(o.dbSync, ob.Name, hospital)
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
func (h *QcDepartmentCtl) Delete() {
	idstr := h.GetString("id")
	if len(idstr) > 0 {
		_, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("failed to parse 'id' value from request, err:", err)
			h.Data["json"] = "invalid 'id' value from request"
			h.ServeJSON()
			return
		}
	}
	err := h.dbSync.DeleteQcObjectWithID(idstr, models.DB_T_DEPARTMENT)
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
func (h *QcDepartmentCtl) Get() {
	hname := h.GetString("hname")
	dname := h.GetString("dname")
	idstr := h.GetString("id")
	var department *models.QcDepartment
	var err error
	if len(hname) > 0 && len(dname) > 0 {
		department, err = h.dbSync.GetQcDepartment(dname, hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get department")
			h.Data["json"] = "invalid id value to get department"
			h.ServeJSON()
			return
		}
		department, err = h.dbSync.GetQcDepartmentWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for department get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = department
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcDepartmentCtl) GetList() {
	pgidx_str := h.GetString("pageidx")
	_, err := strconv.Atoi(pgidx_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pageidx' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pageidx' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	pgsize_str := h.GetString("pagesize")
	_, err = strconv.Atoi(pgsize_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pagesize' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
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
	departments, err := h.dbSync.GetQcDepartmentsCond(pgidx_str, pgsize_str, hid_str)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_DEPARTMENT)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(departments)
	res_data["objects"] = departments
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcDepartmentCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse department id from request")
		h.Data["json"] = "failed to parse department id from request"
		h.ServeJSON()
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		h.logger.LogError("invalid value for department id from request")
		h.Data["json"] = "invalid value department id from request"
		h.ServeJSON()
		return
	}
	department, err := h.dbSync.GetQcDepartmentWithId(id)
	if err != nil {
		h.logger.LogError("failed to get department[", idstr, "] from database, err: ", err)
		h.Data["json"] = "failed to get devmodel[" + idstr + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		department.Name = new_name
	}
	new_hname := h.GetString("hospital")
	if len(new_hname) > 0 {
		new_hospital, err := h.dbSync.GetQcHospital(new_hname)
		if err != nil {
			h.logger.LogError("failed to get hospital[", new_hname, "] info, err: ", err)
			h.Data["json"] = "failed to get hospital[" + new_hname + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		department.Hospital = new_hospital
	}
	err = h.dbSync.UpdateQcDepartment(department)
	if err != nil {
		h.logger.LogError("failed to update department[", idstr, "], err: ", err)
		h.Data["json"] = "failed to update department[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = department
	h.ServeJSON()
	return
}

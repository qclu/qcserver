package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcDevModelCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcDevModelCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcDevModelCtl) Post() {
	mname := o.GetString("methodology")
	mth, err := o.dbSync.GetQcMethodology(mname)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get methodology[" + mname + "] info, err: " + err.Error()
		o.ServeJSON()
	}
	var ob models.QcDevModel
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	o.logger.LogInfo("devmodel: ", ob)
	pob, err := models.CreateQcDevModel(o.dbSync, ob.Name, ob.Model, ob.Release, mth, ob.Annotation)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcDevModelCtl) Delete() {
	hname := h.GetString("name")
	err := h.dbSync.DeleteQcDevModelSQL(hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

// @router / [get]
func (h *QcDevModelCtl) Get() {
	hname := h.GetString("name")
	devmodel, err := h.dbSync.GetQcDevmodel(hname)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = devmodel
	h.ServeJSON()
}

// @router /list [get]
func (h *QcDevModelCtl) GetList() {
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
	devmodels, err := h.dbSync.GetQcDevmodels(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = devmodels
	h.ServeJSON()
}

// @router / [PUT]
func (h *QcDevModelCtl) Update() {
	hname := h.GetString("org_name")
	if len(hname) == 0 {
		h.logger.LogError("failed to parse devmodel name from request")
		h.Data["json"] = "failed to parse devmodel name from request"
		h.ServeJSON()
	}
	devmodel, err := h.dbSync.GetQcDevModel(hname)
	if err != nil {
		h.logger.LogError("failed to get devmodel[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get devmodel[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		devmodel.Name = new_name
	}
	new_release := h.GetString("release")
	if len(new_release) > 0 {
		devmodel.Release = new_release
	}
	new_mth := h.GetString("methodology")
	if len(new_mth) > 0 {
		new_mth_obj, err := h.dbSync.GetQcMethodology(new_mth)
		if err != nil {
			h.logger.LogError("failed to get methodology[", new_mth, "] info, err: ", err)
			h.Data["json"] = "failed to get methodology[" + new_mth + "] info, err: " + err.Error()
			h.ServeJSON()
		}
		devmodel.Methodology = new_mth_obj
	}
	new_model := h.GetString("model")
	if len(new_model) > 0 {
		devmodel.Model = new_model
	}
	err = h.dbSync.UpdateQcDevModel(devmodel)
	if err != nil {
		h.logger.LogError("failed to update devmodel[", hname, "], err: ", err)
		h.Data["json"] = "failed to update devmodel[" + hname + "], err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = devmodel
	h.ServeJSON()
}
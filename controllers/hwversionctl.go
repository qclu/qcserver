package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcHwVersionCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcHwVersionCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcHwVersionCtl) Post() {
	dname := o.GetString("devmodel")
	devmodel, err := o.dbSync.GetQcDevmodel(dname)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get devmodel[" + dname + "] info, err: " + err.Error()
		o.ServeJSON()
	}
	var hwv models.QcHwVersion
	json.Unmarshal(o.Ctx.Input.RequestBody, &hwv)
	pob, err := models.CreateQcHwVersion(o.dbSync, devmodel, hwv.Version, hwv.Annotation)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcHwVersionCtl) Delete() {
	version := h.GetString("version")
	err := h.dbSync.DeleteQcHwVersionSQL(version)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

// @router / [get]
func (h *QcHwVersionCtl) Get() {
	version := h.GetString("version")
	hwv, err := h.dbSync.GetQcHwVersion(version)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = hwv
	h.ServeJSON()
}

// @router /list [get]
func (h *QcHwVersionCtl) GetList() {
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
	hwvs, err := h.dbSync.GetQcHwVersions(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = hwvs
	h.ServeJSON()
}

// @router / [PUT]
func (h *QcHwVersionCtl) Update() {
	version := h.GetString("org_version")
	if len(version) == 0 {
		h.logger.LogError("failed to parse version info from request")
		h.Data["json"] = "failed to parse version info from request"
		h.ServeJSON()
	}
	hwv, err := h.dbSync.GetQcHwVersion(version)
	if err != nil {
		h.logger.LogError("failed to get hwversion[", version, "] from database, err: ", err)
		h.Data["json"] = "failed to get hwversion[" + version + "] from database, err: " + err.Error()
		h.ServeJSON()
	}
	new_version := h.GetString("version")
	if len(new_version) > 0 {
		hwv.Version = new_version
	}
	new_anno := h.GetString("annotation")
	if len(new_anno) > 0 {
		hwv.Annotation = new_anno
	}
	new_devmodel := h.GetString("devmodel")
	if len(new_devmodel) > 0 {
		new_devmodel_obj, err := h.dbSync.GetQcDevmodel(new_devmodel)
		if err != nil {
			h.logger.LogError("failed to get devmodel[", new_devmodel, "] info, err: ", err)
			h.Data["json"] = "failed to get devmodel[" + new_devmodel + "] info, err: " + err.Error()
			h.ServeJSON()
		}
		hwv.DevModel = new_devmodel_obj
	}
	err = h.dbSync.UpdateQcHwVersion(hwv)
	if err != nil {
		h.logger.LogError("failed to update hwversion[", version, "], err: ", err)
		h.Data["json"] = "failed to update hwversion[" + version + "], err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = hwv
	h.ServeJSON()
}

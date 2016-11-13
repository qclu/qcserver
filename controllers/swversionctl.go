package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcSwVersionCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcSwVersionCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcSwVersionCtl) Post() {
	hwversion := o.GetString("hwversion")
	hwv, err := o.dbSync.GetQcHwVersion(hwversion)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get hwversion[" + hwversion + "] info, err: " + err.Error()
		o.ServeJSON()
	}
	var swv models.QcSwVersion
	json.Unmarshal(o.Ctx.Input.RequestBody, &swv)
	pob, err := models.CreateQcSwVersion(o.dbSync, hwv, swv.Version, swv.SwType, swv.Description)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
	}
	o.Data["json"] = pob
	o.ServeJSON()
}

// @router / [delete]
func (h *QcSwVersionCtl) Delete() {
	version := h.GetString("version")
	err := h.dbSync.DeleteQcSwVersionSQL(version)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
}

// @router / [get]
func (h *QcSwVersionCtl) Get() {
	version := h.GetString("version")
	Swv, err := h.dbSync.GetQcSwVersion(version)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = Swv
	h.ServeJSON()
}

// @router /list [get]
func (h *QcSwVersionCtl) GetList() {
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
	swvs, err := h.dbSync.GetQcSwVersions(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = swvs
	h.ServeJSON()
}

// @router / [PUT]
func (h *QcSwVersionCtl) Update() {
	version := h.GetString("org_version")
	if len(version) == 0 {
		h.logger.LogError("failed to parse version info from request")
		h.Data["json"] = "failed to parse version info from request"
		h.ServeJSON()
	}
	swv, err := h.dbSync.GetQcSwVersion(version)
	if err != nil {
		h.logger.LogError("failed to get swversion[", version, "] from database, err: ", err)
		h.Data["json"] = "failed to get swversion[" + version + "] from database, err: " + err.Error()
		h.ServeJSON()
	}
	new_version := h.GetString("version")
	if len(new_version) > 0 {
		swv.Version = new_version
	}
	new_type := h.GetString("swtype")
	if len(new_type) > 0 {
		swv.SwType = new_type
	}
	new_description := h.GetString("description")
	if len(new_description) > 0 {
		swv.Description = new_description
	}
	new_hwv := h.GetString("hwversion")
	if len(new_hwv) > 0 {
		new_hwv_obj, err := h.dbSync.GetQcHwVersion(new_hwv)
		if err != nil {
			h.logger.LogError("failed to get hwversion[", new_hwv, "] info, err: ", err)
			h.Data["json"] = "failed to get hwversion[" + new_hwv + "] info, err: " + err.Error()
			h.ServeJSON()
		}
		swv.HwVersion = new_hwv_obj
	}
	err = h.dbSync.UpdateQcSwVersion(swv)
	if err != nil {
		h.logger.LogError("failed to update swversion[", version, "], err: ", err)
		h.Data["json"] = "failed to update swversion[" + version + "], err: " + err.Error()
		h.ServeJSON()
	}
	h.Data["json"] = swv
	h.ServeJSON()
}
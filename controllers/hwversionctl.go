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
		return
	}
	var hwv models.QcHwVersion
	json.Unmarshal(o.Ctx.Input.RequestBody, &hwv)
	pob, err := models.CreateQcHwVersion(o.dbSync, devmodel, hwv.Version, hwv.Annotation)
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
func (h *QcHwVersionCtl) Delete() {
	version := h.GetString("version")
	err := h.dbSync.DeleteQcHwVersionSQL(version)
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
func (h *QcHwVersionCtl) Get() {
	version := h.GetString("version")
	idstr := h.GetString("id")
	var hwv *models.QcHwVersion
	var err error
	if len(version) > 0 {
		hwv, err = h.dbSync.GetQcHwVersion(version)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get hardware version")
			h.Data["json"] = "invalid id value to get hardware version"
			h.ServeJSON()
			return
		}
		hwv, err = h.dbSync.GetQcHwVersionWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for hardware version get"
		h.ServeJSON()
		return
	}

	h.Data["json"] = hwv
	h.ServeJSON()
	return
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
	hwvs, err := h.dbSync.GetQcHwVersions(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_HWVERSION)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = entcnt
	res_data["objects"] = hwvs
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcHwVersionCtl) Update() {
	version := h.GetString("org_version")
	if len(version) == 0 {
		h.logger.LogError("failed to parse version info from request")
		h.Data["json"] = "failed to parse version info from request"
		h.ServeJSON()
		return
	}
	hwv, err := h.dbSync.GetQcHwVersion(version)
	if err != nil {
		h.logger.LogError("failed to get hwversion[", version, "] from database, err: ", err)
		h.Data["json"] = "failed to get hwversion[" + version + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
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
			return
		}
		hwv.DevModel = new_devmodel_obj
	}
	err = h.dbSync.UpdateQcHwVersion(hwv)
	if err != nil {
		h.logger.LogError("failed to update hwversion[", version, "], err: ", err)
		h.Data["json"] = "failed to update hwversion[" + version + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = hwv
	h.ServeJSON()
	return
}

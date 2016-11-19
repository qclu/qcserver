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
		return
	}
	var swv models.QcSwVersion
	json.Unmarshal(o.Ctx.Input.RequestBody, &swv)
	pob, err := models.CreateQcSwVersion(o.dbSync, hwv, swv.Version, swv.SwType, swv.Description)
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
func (h *QcSwVersionCtl) Delete() {
	id := h.GetString("id")
	err := h.dbSync.DeleteQcObjectSQL(id, models.DB_T_SWVERSION)
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
func (h *QcSwVersionCtl) Get() {
	version := h.GetString("version")
	idstr := h.GetString("id")
	var swv *models.QcSwVersion
	var err error
	if len(version) > 0 {
		swv, err = h.dbSync.GetQcSwVersion(version)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get software version")
			h.Data["json"] = "invalid id value to get software version"
			h.ServeJSON()
			return
		}
		swv, err = h.dbSync.GetQcSwVersionWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for software version get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = swv
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcSwVersionCtl) GetList() {
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
	h.logger.LogInfo("list hospital info(pageidx: ", pgidx_str, ", pagesize: ", pgsize_str, ")")
	devid_str := h.GetString("devid")
	if len(devid_str) > 0 {
		_, err = strconv.Atoi(devid_str)
		if err != nil {
			h.logger.LogError("invalid value for param 'devid' from request, err: ", err)
			h.Data["json"] = "invalid value for param 'devid' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}

	hwv_str := h.GetString("hwvid")
	if len(hwv_str) > 0 {
		_, err = strconv.Atoi(hwv_str)
		if err != nil {
			h.logger.LogError("invalid value for param 'hwvid' from request, err: ", err)
			h.Data["json"] = "invalid value for param 'hwvid' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}

	version := h.GetString("version")
	swvs, err := h.dbSync.GetQcSwVersionsCond(pgidx_str, pgsize_str, devid_str, hwv_str, version)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_SWVERSION)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(swvs)
	res_data["objects"] = swvs
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcSwVersionCtl) Update() {
	version := h.GetString("org_version")
	if len(version) == 0 {
		h.logger.LogError("failed to parse version info from request")
		h.Data["json"] = "failed to parse version info from request"
		h.ServeJSON()
		return
	}
	swv, err := h.dbSync.GetQcSwVersion(version)
	if err != nil {
		h.logger.LogError("failed to get swversion[", version, "] from database, err: ", err)
		h.Data["json"] = "failed to get swversion[" + version + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
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
			return
		}
		swv.HwVersion = new_hwv_obj
	}
	err = h.dbSync.UpdateQcSwVersion(swv)
	if err != nil {
		h.logger.LogError("failed to update swversion[", version, "], err: ", err)
		h.Data["json"] = "failed to update swversion[" + version + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = swv
	h.ServeJSON()
	return
}

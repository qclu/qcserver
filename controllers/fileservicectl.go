package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcFileServiceCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcFileServiceCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcFileServiceCtl) Post() {
	idstr := o.GetString("swvid")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		o.logger.LogError("failed to parse 'swvid' from request, err: ", err)
		o.Data["json"] = "invalid parse 'swvid' from request, err: " + err.Error()
		o.ServeJSON()
		return
	}
	swv, err := o.dbSync.GetQcSwVersionWithId(id)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get sw version[" + idstr + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	filename := "/data/" + strconv.FormatInt(swv.Id, 10) + "_" + swv.Version + swv.SwType + ".zip"
	err = o.SaveToFile("file", filename)
	if err != nil {
		o.logger.LogError("failed to save file: ", filename, " err: ", err)
		o.Data["json"] = string("failed to save file: ") + filename + string(" err: ") + err.Error()
		o.ServeJSON()
		return
	}
	o.Data["json"] = "OK"
	o.ServeJSON()
	return
}

// @router / [get]
func (h *QcFileServiceCtl) Get() {
	idstr := h.GetString("swvid")
	var swv *models.QcSwVersion
	if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get sw version")
			h.Data["json"] = "invalid id value to get sw version"
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
		h.Data["json"] = "invalid parameter for file service"
		h.ServeJSON()
		return
	}
	filename := strconv.FormatInt(swv.Id, 10) + "_" + swv.Version + swv.SwType + ".zip"
	h.Data["json"] = "/swpackage/" + filename
	h.ServeJSON()
	return
}

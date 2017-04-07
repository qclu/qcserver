package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcLogTypeCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcLogTypeCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcLogTypeCtl) Post() {
	var ob models.QcLogType
	var err error
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	o.logger.LogInfo("RequestBody: ", string(o.Ctx.Input.RequestBody))
	o.logger.LogInfo("Logtype create request: ", ob)
	if o.dbSync == nil {
		o.logger.LogError("dbSync is uninitialized when trying to create QcLogType, request: ", o.Ctx)
		o.Abort("501")
	}
	pob, err := models.CreateQcLogType(o.dbSync, ob.Type, ob.Level, ob.Content)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = string("database operation err:") + err.Error()
		o.ServeJSON()
		return
	}
	o.Data["json"] = map[string]models.QcLogType{"Object": *pob}
	o.ServeJSON()
	return
}

// @router / [delete]
func (h *QcLogTypeCtl) Delete() {
	idstr := h.GetString("id")
	_, err := strconv.Atoi(idstr)
	if err != nil {
		h.logger.LogError("invalid id value to delete logtype")
		h.Data["json"] = "invalid id value to delete log type"
		h.ServeJSON()
		return
	}
	err = h.dbSync.DeleteQcObjectWithID(idstr, models.DB_T_LOGTYPE)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "internal database operation error"
		h.ServeJSON()
		return
	}
	h.Data["json"] = "delete success"
	h.ServeJSON()
	return
}

// @router / [get]
func (h *QcLogTypeCtl) Get() {
	idstr := h.GetString("id")
	var logtype *models.QcLogType
	if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get logtype")
			h.Data["json"] = "invalid id value to get logtype"
			h.ServeJSON()
			return
		}
		logtype, err = h.dbSync.GetQcLogTypeWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for logtype get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = logtype
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcLogTypeCtl) GetList() {
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
	h.logger.LogInfo("list logtype info(pageidx: ", pgidx, ", pagesize: ", pgsize, ")")
	idstr := h.GetString("id")
	if len(idstr) > 0 {
		h.logger.LogInfo("log id: ", idstr)
		_, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("failed to parse 'id' from request, err: ", err)
			h.Data["json"] = "invalid parse 'id' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	lvlstr := h.GetString("level")
	if len(lvlstr) > 0 {
		h.logger.LogInfo("log level: ", lvlstr)
		_, err := strconv.Atoi(lvlstr)
		if err != nil {
			h.logger.LogError("failed to parse 'level' from request, err: ", err)
			h.Data["json"] = "invalid parse 'level' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	logtypes, err := h.dbSync.GetQcLogTypeCond(pgidx_str, pgsize_str, idstr, lvlstr)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Abort("501")
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_HOSPITAL)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(logtypes)
	res_data["objects"] = logtypes
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcLogTypeCtl) Update() {
	//idstr := h.GetString("id")
	//if len(idstr) == 0 {
	//	h.logger.LogError("failed to parse logtype id from request")
	//	h.Data["json"] = "failed to parse logtype id from request"
	//	h.ServeJSON()
	//	return
	//}
	//id, _ := strconv.Atoi(idstr)
	//_, err := h.dbSync.GetQcLogTypeWithId(id)
	//if err != nil {
	//	h.logger.LogError("failed to get logtype[", id, "] from database, err: ", err)
	//	h.Data["json"] = "failed to get logtype[" + idstr + "] from database, err: " + err.Error()
	//	h.ServeJSON()
	//	return
	//}
	var newdata models.QcLogType
	var err error
	err = json.Unmarshal(h.Ctx.Input.RequestBody, &newdata)
	if err != nil {
		h.logger.LogError("failed to unmarshal new logtype data, err: ", err)
		h.Data["json"] = "failed to unmarshal new logtype data, err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.logger.LogInfo("data to update: ", newdata)
	err = h.dbSync.UpdateQcLogType(&newdata)
	if err != nil {
		h.logger.LogError("failed to update logtype[", newdata, "], err: ", err)
		h.Data["json"] = "failed to update logtype[" + string(h.Ctx.Input.RequestBody) + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = newdata
	h.ServeJSON()
	return
}

package controllers

import (
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcDevLogCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcDevLogCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router /list [get]
func (h *QcDevLogCtl) GetList() {
	pgidx_str := h.GetString("pageidx")
	if len(pgidx_str) > 0 {
		h.logger.LogInfo("pageidx: ", pgidx_str)
		_, err := strconv.Atoi(pgidx_str)
		if err != nil {
			h.logger.LogError("failed to parse 'pageidx' from request, err: ", err)
			h.Data["json"] = "invalid parse 'pageidx' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	pgsize_str := h.GetString("pagesize")
	if len(pgsize_str) > 0 {
		h.logger.LogInfo("pagesize: ", pgsize_str)
		_, err := strconv.Atoi(pgsize_str)
		if err != nil {
			h.logger.LogError("failed to parse 'pagesize' from request, err: ", err)
			h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	msgtype_str := h.GetString("type")
	if len(msgtype_str) > 0 {
		h.logger.LogInfo("msgtype: ", msgtype_str)
		_, err := strconv.Atoi(msgtype_str)
		if err != nil {
			h.logger.LogError("failed to parse 'type' from request, err: ", err)
			h.Data["json"] = "invalid parse 'type' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	devsn_str := h.GetString("sn")
	if len(devsn_str) > 0 {
		h.logger.LogInfo("devsn: ", devsn_str)
		_, err := strconv.Atoi(devsn_str)
		if err != nil {
			h.logger.LogError("failed to parse 'sn' from request, err: ", err)
			h.Data["json"] = "invalid parse 'sn' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	logs, err := h.dbSync.GetQcDevLogCond(pgidx_str, pgsize_str, msgtype_str, devsn_str)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get dev logs, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_QCPRODUCT)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(logs)
	res_data["objects"] = logs
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

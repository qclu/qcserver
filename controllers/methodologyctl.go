package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcMethodologyCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcMethodologyCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcMethodologyCtl) Post() {
	var ob models.QcMethodology
	var err error
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcMethodology(o.dbSync, ob.Name, ob.Annotation)
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
func (h *QcMethodologyCtl) Delete() {
	hname := h.GetString("name")
	err := h.dbSync.DeleteQcMethdologySQL(hname)
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
func (h *QcMethodologyCtl) Get() {
	hname := h.GetString("name")
	idstr := h.GetString("id")
	var mt *models.QcMethodology
	var err error
	if len(hname) > 0 {
		mt, err = h.dbSync.GetQcMethodology(hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get methodology")
			h.Data["json"] = "invalid id value to get methodology"
			h.ServeJSON()
			return
		}
		mt, err = h.dbSync.GetQcMethodologyWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for methodology get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = mt
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcMethodologyCtl) GetList() {
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
	ms, err := h.dbSync.GetQcMethodologys(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_METHODOLOGY)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(ms)
	res_data["objects"] = ms
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcMethodologyCtl) Update() {
	hname := h.GetString("org_name")
	if len(hname) == 0 {
		h.logger.LogError("failed to parse admin name from request")
		h.Data["json"] = "failed to parse admin name from request"
		h.ServeJSON()
		return
	}
	mth, err := h.dbSync.GetQcMethodology(hname)
	if err != nil {
		h.logger.LogError("failed to get methodology[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get methodology[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		mth.Name = new_name
	}
	new_anno := h.GetString("anno")
	if len(new_anno) > 0 {
		mth.Annotation = new_anno
	}
	err = h.dbSync.UpdateQcMethodology(mth)
	if err != nil {
		h.logger.LogError("failed to update methodology[", hname, "], err: ", err)
		h.Data["json"] = "failed to update methodology[" + hname + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = mth
	h.ServeJSON()
	return
}

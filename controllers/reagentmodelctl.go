package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcReagentModelCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcReagentModelCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcReagentModelCtl) Post() {
	devmodel_str := o.GetString("devmodel")
	devmodel, err := o.dbSync.GetQcDevmodel(devmodel_str)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get devmodel[" + devmodel_str + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	var ob models.QcReagentModel
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcReagentModel(o.dbSync, ob.Name, ob.Annotation, devmodel)
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
func (h *QcReagentModelCtl) Delete() {
	idstr := h.GetString("id")
	err := h.dbSync.DeleteQcObjectSQL(idstr, models.DB_T_REGMODEL)
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
func (h *QcReagentModelCtl) Get() {
	hname := h.GetString("name")
	idstr := h.GetString("id")
	var regmodel *models.QcReagentModel
	var err error
	if len(hname) > 0 {
		regmodel, err = h.dbSync.GetQcReagentModel(hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get reagent model")
			h.Data["json"] = "invalid id value to get reagent model"
			h.ServeJSON()
			return
		}
		regmodel, err = h.dbSync.GetQcReagentModelWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for reagent model get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = regmodel
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcReagentModelCtl) GetList() {
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
	h.logger.LogInfo("list regmodel info(pageidx: ", pgidx, ", pagesize: ", pgsize, ")")
	devmodelid_str := h.GetString("devmodelid")
	if len(devmodelid_str) > 0 {
		_, err = strconv.Atoi(devmodelid_str)
		if err != nil {
			h.logger.LogError("failed to parse 'devmodelid' from request, err: ", err)
			h.Data["json"] = "invalid value for 'devmodel' from request, err: " + err.Error()
			h.ServeJSON()
			return
		}
	}
	name := h.GetString("name")
	regmodels, err := h.dbSync.GetQcReagentModelsCond(pgidx_str, pgsize_str, devmodelid_str, name)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get reagent models, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_REGMODEL)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(regmodels)
	res_data["objects"] = regmodels
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcReagentModelCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse regmodel id from request")
		h.Data["json"] = "failed to parse regmodel id from request"
		h.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(idstr)
	regmodel, err := h.dbSync.GetQcReagentModelWithId(id)
	if err != nil {
		h.logger.LogError("failed to get reagent model[", id, "] from database, err: ", err)
		h.Data["json"] = "failed to get reagent model[" + idstr + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		regmodel.Name = new_name
	}
	new_devmodel := h.GetString("devmodel")
	if len(new_devmodel) > 0 {
		new_devmodel_obj, err := h.dbSync.GetQcDevModel(new_devmodel)
		if err != nil {
			h.logger.LogError("failed to get devmodel[", new_devmodel, "] info, err: ", err)
			h.Data["json"] = "failed to get devmodel[" + new_devmodel + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		regmodel.DevModel = new_devmodel_obj
	}
	new_anno := h.GetString("annotation")
	if len(new_anno) > 0 {
		regmodel.Annotation = new_anno
	}
	err = h.dbSync.UpdateQcReagentModel(regmodel)
	if err != nil {
		h.logger.LogError("failed to update regmodel[", id, "], err: ", err)
		h.Data["json"] = "failed to update regmodel[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = regmodel
	h.ServeJSON()
	return
}

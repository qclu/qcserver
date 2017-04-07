package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcReagentProduceCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcReagentProduceCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcReagentProduceCtl) Post() {
	regmodelid_str := o.GetString("regmodel")
	regmodelid, err := strconv.Atoi(regmodelid_str)
	if err != nil {
		o.logger.LogError("invalid id value to get reagent model")
		o.Data["json"] = "invalid id value to get reagent model"
		o.ServeJSON()
		return
	}
	regmodel, err := o.dbSync.GetQcReagentModelWithId(regmodelid)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get reagentmodel[" + regmodelid_str + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	var ob models.QcReagentProduce
	err = json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	if err != nil {
		o.logger.LogError("Failed to unmarshal Reagent Produce info, err: ", err)
		o.Data["json"] = string("Failed to unmarshal Reagent Produce info, err: ") + err.Error()
		o.ServeJSON()
		return
	}
	o.logger.LogInfo("Regproduce info: ", ob)
	pob, err := models.CreateQcReagentProduce(o.dbSync, ob.SerialNum, ob.ExpiredTime, ob.Annotation, regmodel)
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
func (h *QcReagentProduceCtl) Delete() {
	idstr := h.GetString("id")
	err := h.dbSync.DeleteQcObjectWithID(idstr, models.DB_T_REGPRODUCE)
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
func (h *QcReagentProduceCtl) Get() {
	serial := h.GetString("serialnum")
	idstr := h.GetString("id")
	var regproduce *models.QcReagentProduce
	var err error
	if len(serial) > 0 {
		regproduce, err = h.dbSync.GetQcReagentProduce(serial)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get reagent produce")
			h.Data["json"] = "invalid id value to get reagent produce"
			h.ServeJSON()
			return
		}
		regproduce, err = h.dbSync.GetQcReagentProduceWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for reagent produce get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = regproduce
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcReagentProduceCtl) GetList() {
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
	reg_str := h.GetString("reagentid")
	h.logger.LogInfo("reagent: ", reg_str)
	var condition string
	if len(reg_str) > 0 {
		condition = "AND reg_model_id=" + reg_str
	}
	h.logger.LogInfo("list reagentproduce info(pageidx: ", pgidx, ", pagesize: ", pgsize, ")")
	regproduces, err := h.dbSync.GetQcReagentProduces(pgidx, pgsize, condition)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get reagent produces, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_REGPRODUCE)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(regproduces)
	res_data["objects"] = regproduces
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcReagentProduceCtl) Update() {
	idstr := h.GetString("id")
	if len(idstr) == 0 {
		h.logger.LogError("failed to parse regproduce id from request")
		h.Data["json"] = "failed to parse regproduce id from request"
		h.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(idstr)
	regproduce, err := h.dbSync.GetQcReagentProduceWithId(id)
	if err != nil {
		h.logger.LogError("failed to get reagent produce[", id, "] from database, err: ", err)
		h.Data["json"] = "failed to get reagent produce[" + idstr + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_serial := h.GetString("serialnum")
	if len(new_serial) > 0 {
		regproduce.SerialNum = new_serial
		h.logger.LogInfo("serial updata to ", new_serial)
	}
	new_expiredtime := h.GetString("expiredtime")
	if len(new_expiredtime) > 0 {
		regproduce.ExpiredTime = new_expiredtime
	}
	new_regmodel := h.GetString("regmodel")
	if len(new_regmodel) > 0 {
		new_regmodel_obj, err := h.dbSync.GetQcReagentModel(new_regmodel)
		if err != nil {
			h.logger.LogError("failed to get reagent model[", new_regmodel, "] info, err: ", err)
			h.Data["json"] = "failed to get reagent model[" + new_regmodel + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		regproduce.RegModel = new_regmodel_obj
	}
	new_anno := h.GetString("annotation")
	if len(new_anno) > 0 {
		regproduce.Annotation = new_anno
	}
	err = h.dbSync.UpdateQcReagentProduce(regproduce)
	if err != nil {
		h.logger.LogError("failed to update reagent produce[", id, "], err: ", err)
		h.Data["json"] = "failed to update reagent model[" + idstr + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = regproduce
	h.ServeJSON()
	return
}

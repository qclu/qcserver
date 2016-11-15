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
	regmodel_str := o.GetString("reagentmodel")
	regmodel, err := o.dbSync.GetQcReagentModel(regmodel_str)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get reagentmodel[" + regmodel_str + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	var ob models.QcReagentProduce
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcReagentProduce(o.dbSync, ob.SerialNum, ob.LotNum, ob.ExpiredTime, ob.Annotation, regmodel)
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
	hname := h.GetString("serialnum")
	err := h.dbSync.DeleteQcReagentProduceSQL(hname)
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
	regproduce, err := h.dbSync.GetQcReagentProduce(serial)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "database operation err: " + err.Error()
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
	h.logger.LogInfo("list hospital info(pageidx: ", pgidx, ", pagesize: ", pgsize, ")")
	regproduces, err := h.dbSync.GetQcReagentProduces(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get reagent produces, err: " + err.Error()
		h.ServeJSON()
		return
	}
	entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_REGPRODUCE)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = entcnt
	res_data["objects"] = regproduces
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcReagentProduceCtl) Update() {
	serial := h.GetString("org_serialnum")
	if len(serial) == 0 {
		h.logger.LogError("failed to parse regproduce serial from request")
		h.Data["json"] = "failed to parse regproduce serial from request"
		h.ServeJSON()
		return
	}
	regproduce, err := h.dbSync.GetQcReagentProduce(serial)
	if err != nil {
		h.logger.LogError("failed to get reagent produce[", serial, "] from database, err: ", err)
		h.Data["json"] = "failed to get reagent produce[" + serial + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_serial := h.GetString("serialnum")
	if len(new_serial) > 0 {
		regproduce.SerialNum = new_serial
	}
	new_lotnum := h.GetString("lotnum")
	if len(new_lotnum) > 0 {
		regproduce.LotNum = new_lotnum
	}
	new_regmodel := h.GetString("regmodl")
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
		h.logger.LogError("failed to update reagent produce[", serial, "], err: ", err)
		h.Data["json"] = "failed to update reagent model[" + serial + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = regproduce
	h.ServeJSON()
	return
}
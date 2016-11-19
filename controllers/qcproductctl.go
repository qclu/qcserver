package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcQcProductCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcQcProductCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcQcProductCtl) Post() {
	regmodel_str := o.GetString("reagentmodel")
	regmodel, err := o.dbSync.GetQcReagentModel(regmodel_str)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get regmodel[" + regmodel_str + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	var ob models.QcQcProduct
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	pob, err := models.CreateQcQcProduct(o.dbSync, regmodel, ob.Tea, ob.Cv, ob.Percent, ob.FixedDeviation, ob.Nsd, ob.Range, ob.Name, ob.Annotation)
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
func (h *QcQcProductCtl) Delete() {
	idstr := h.GetString("id")
	err := h.dbSync.DeleteQcObjectSQL(idstr, models.DB_T_QCPRODUCT)
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
func (h *QcQcProductCtl) Get() {
	hname := h.GetString("name")
	idstr := h.GetString("id")
	var qcp *models.QcQcProduct
	var err error
	if len(hname) > 0 {
		qcp, err = h.dbSync.GetQcQcProduct(hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get qcproduct")
			h.Data["json"] = "invalid id value to get qcproduct"
			h.ServeJSON()
			return
		}
		qcp, err = h.dbSync.GetQcQcProductWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for qcproduct get"
		h.ServeJSON()
		return
	}
	h.Data["json"] = qcp
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcQcProductCtl) GetList() {
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
	qcps, err := h.dbSync.GetQcQcProducts(pgidx, pgsize, "")
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "failed to get qcproducts, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_QCPRODUCT)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(qcps)
	res_data["objects"] = qcps
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcQcProductCtl) Update() {
	hname := h.GetString("org_name")
	if len(hname) == 0 {
		h.logger.LogError("failed to parse qcproduct name from request")
		h.Data["json"] = "failed to parse qcproduct name from request"
		h.ServeJSON()
		return
	}
	qcp, err := h.dbSync.GetQcQcProduct(hname)
	if err != nil {
		h.logger.LogError("failed to get qcproduct[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get qcproduct[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		qcp.Name = new_name
	}
	new_regmodel := h.GetString("reagentmodel")
	if len(new_regmodel) > 0 {
		new_regmodel_obj, err := h.dbSync.GetQcReagentModel(new_regmodel)
		if err != nil {
			h.logger.LogError("failed to get regmodel[", new_regmodel, "] info, err: ", err)
			h.Data["json"] = "failed to get regmodel[" + new_regmodel + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		qcp.RegModel = new_regmodel_obj
	}
	new_anno := h.GetString("annotation")
	if len(new_anno) > 0 {
		qcp.Annotation = new_anno
	}
	new_tea := h.GetString("tea")
	if len(new_tea) > 0 {
		tea, err := strconv.ParseFloat(new_tea, 64)
		if err != nil {
			h.logger.LogError("invalid parameter for 'tea'[", new_tea, "] info, err: ", err)
			h.Data["json"] = "invalid parameter for 'tea'[" + new_tea + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		qcp.Tea = tea
	}
	new_rg := h.GetString("range")
	if len(new_rg) > 0 {
		qcp.Range = new_rg
	}
	new_nsd := h.GetString("nsd")
	if len(new_nsd) > 0 {
		nsd, err := strconv.ParseFloat(new_nsd, 64)
		if err != nil {
			h.logger.LogError("invalid parameter for 'nsd'[", new_nsd, "] info, err: ", err)
			h.Data["json"] = "invalid parameter for 'nsd'[" + new_nsd + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		qcp.Nsd = nsd
	}
	new_fixd := h.GetString("fixdeviation")
	if len(new_fixd) > 0 {
		fixd, err := strconv.ParseFloat(new_fixd, 64)
		if err != nil {
			h.logger.LogError("invalid parameter for 'fixdeviation'[", new_fixd, "] info, err: ", err)
			h.Data["json"] = "invalid parameter for 'fixdeviation'[" + new_fixd + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		qcp.FixedDeviation = fixd
	}
	new_percent := h.GetString("percent")
	if len(new_percent) > 0 {
		percent, err := strconv.ParseFloat(new_percent, 64)
		if err != nil {
			h.logger.LogError("invalid parameter for 'percent'[", new_percent, "] info, err: ", err)
			h.Data["json"] = "invalid parameter for 'percent'[" + new_percent + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		qcp.Percent = percent
	}
	new_cv := h.GetString("cv")
	if len(new_cv) > 0 {
		cv, err := strconv.ParseFloat(new_cv, 64)
		if err != nil {
			h.logger.LogError("invalid parameter for 'cv'[", new_cv, "] info, err: ", err)
			h.Data["json"] = "invalid parameter for 'cv'[" + new_cv + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		qcp.Cv = cv
	}
	err = h.dbSync.UpdateQcQcProduct(qcp)
	if err != nil {
		h.logger.LogError("failed to update qcproduct[", hname, "], err: ", err)
		h.Data["json"] = "failed to update qcproduct[" + hname + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = qcp
	h.ServeJSON()
	return
}

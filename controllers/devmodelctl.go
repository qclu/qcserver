package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"qcserver/util/log"
	"strconv"
)

// Operations about object
type QcDevModelCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcDevModelCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [post]
func (o *QcDevModelCtl) Post() {
	mname := o.GetString("methodology")
	mth, err := o.dbSync.GetQcMethodology(mname)
	if err != nil {
		o.logger.LogError("database operation err: ", err)
		o.Data["json"] = "failed to get methodology[" + mname + "] info, err: " + err.Error()
		o.ServeJSON()
		return
	}
	var ob models.QcDevModel
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	o.logger.LogInfo("devmodel: ", ob)
	pob, err := models.CreateQcDevModel(o.dbSync, ob.Name, ob.Model, ob.Release, mth, ob.Annotation)
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
func (h *QcDevModelCtl) Delete() {
	hname := h.GetString("name")
	err := h.dbSync.DeleteQcDevModelSQL(hname)
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
func (h *QcDevModelCtl) Get() {
	hname := h.GetString("name")
	idstr := h.GetString("id")
	var devmodel *models.QcDevModel
	var err error
	if len(hname) > 0 {
		devmodel, err = h.dbSync.GetQcDevmodel(hname)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else if len(idstr) > 0 {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			h.logger.LogError("invalid id value to get devmodel")
			h.Data["json"] = "invalid id value to get devmodel"
			h.ServeJSON()
			return
		}
		devmodel, err = h.dbSync.GetQcDevmodelWithId(id)
		if err != nil {
			h.logger.LogError("database operation err: ", err)
			h.Data["json"] = "database operation err: " + err.Error()
			h.ServeJSON()
			return
		}
	} else {
		h.Data["json"] = "invalid parameter for devmodel get"
		h.ServeJSON()
		return

	}
	h.Data["json"] = devmodel
	h.ServeJSON()
	return
}

// @router /list [get]
func (h *QcDevModelCtl) GetList() {
	pgidx_str := h.GetString("pageidx")
	pgidx, err := strconv.Atoi(pgidx_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pageidx' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pageidx' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	pgsize_str := h.GetString("pagesize")
	pgsize, err := strconv.Atoi(pgsize_str)
	if err != nil {
		h.logger.LogError("failed to parse 'pagesize' from request, err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	name_str := h.GetString("name")
	h.logger.LogInfo("name:", name_str)
	var condition string
	if len(name_str) > 0 {
		condition = " and name='" + name_str + "'"
	}
	model_str := h.GetString("model")
	if len(model_str) > 0 {
		if len(condition) > 0 {
			condition = condition + " and " + "model='" + model_str + "'"
		} else {
			condition = " and model='" + model_str + "'"
		}
	}
	h.logger.LogInfo("condition:", condition)
	devmodels, err := h.dbSync.GetQcDevmodels(pgidx, pgsize, condition)
	if err != nil {
		h.logger.LogError("database operation err: ", err)
		h.Data["json"] = "invalid parse 'pagesize' from request, err: " + err.Error()
		h.ServeJSON()
		return
	}
	//entcnt, _ := h.dbSync.GetTotalCnt(models.DB_T_DEVMODEL)
	res_data := make(map[string]interface{})
	res_data["totalnum"] = len(devmodels)
	res_data["objects"] = devmodels
	h.Data["json"] = res_data
	h.ServeJSON()
	return
}

// @router / [PUT]
func (h *QcDevModelCtl) Update() {
	hname := h.GetString("org_name")
	if len(hname) == 0 {
		h.logger.LogError("failed to parse devmodel name from request")
		h.Data["json"] = "failed to parse devmodel name from request"
		h.ServeJSON()
		return
	}
	devmodel, err := h.dbSync.GetQcDevModel(hname)
	if err != nil {
		h.logger.LogError("failed to get devmodel[", hname, "] from database, err: ", err)
		h.Data["json"] = "failed to get devmodel[" + hname + "] from database, err: " + err.Error()
		h.ServeJSON()
		return
	}
	new_name := h.GetString("name")
	if len(new_name) > 0 {
		devmodel.Name = new_name
	}
	new_release := h.GetString("release")
	if len(new_release) > 0 {
		devmodel.Release = new_release
	}
	new_mth := h.GetString("methodology")
	if len(new_mth) > 0 {
		new_mth_obj, err := h.dbSync.GetQcMethodology(new_mth)
		if err != nil {
			h.logger.LogError("failed to get methodology[", new_mth, "] info, err: ", err)
			h.Data["json"] = "failed to get methodology[" + new_mth + "] info, err: " + err.Error()
			h.ServeJSON()
			return
		}
		devmodel.Methodology = new_mth_obj
	}
	new_model := h.GetString("model")
	if len(new_model) > 0 {
		devmodel.Model = new_model
	}
	err = h.dbSync.UpdateQcDevModel(devmodel)
	if err != nil {
		h.logger.LogError("failed to update devmodel[", hname, "], err: ", err)
		h.Data["json"] = "failed to update devmodel[" + hname + "], err: " + err.Error()
		h.ServeJSON()
		return
	}
	h.Data["json"] = devmodel
	h.ServeJSON()
	return
}

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/util/log"
)

// Operations about object
type QcMapGeneratorCtl struct {
	beego.Controller
	logger *log.Log
}

func (this *QcMapGeneratorCtl) Prepare() {
	this.logger = log.GetLog()
}

type QcGisInfo struct {
	Longitude float64
	Latitude  float64
}

type QcMapLocationInfo struct {
	City string
	Gis  []QcGisInfo
}

// @router / [post]
func (h *QcMapGeneratorCtl) Post() {
	var loc_param QcMapLocationInfo
	h.logger.LogInfo("request body: ", string(h.Ctx.Input.RequestBody))
	err := json.Unmarshal(h.Ctx.Input.RequestBody, &loc_param)
	if err != nil {
		h.logger.LogError("error: ", err.Error())
	}
	h.logger.LogInfo("location param: ", loc_param.City)
	h.Data["City"] = loc_param.City
	h.Data["Gis"] = loc_param.Gis
	h.TplName = "bdmaptest.tpl"
}

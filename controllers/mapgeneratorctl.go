package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	"qcserver/models"
	"strconv"
	"qcserver/util/log"
)

// Operations about object
type QcMapGeneratorCtl struct {
	beego.Controller
	logger *log.Log
	dbSync *models.DBSync
}

func (this *QcMapGeneratorCtl) Prepare() {
	this.logger = log.GetLog()
	this.dbSync = models.GetDBSync()
}

// @router / [get]
func (h *QcMapGeneratorCtl) Get() {
	h.logger.LogInfo("request body: ", string(h.Ctx.Input.RequestBody))
	cityname := h.GetString("city")
	tmpgisdata, err := h.dbSync.GetAllGisInfo()
	if err != nil {
		h.logger.LogError("error: ", err.Error())
	}
	type GisInfo struct {
	        Longitude float64
        	Latitude  float64
	}
	var gislist []GisInfo
	for index := 0; index < len(tmpgisdata); index++  {
		var gisdata GisInfo
		gisdata.Longitude, _ = strconv.ParseFloat(tmpgisdata[index].Longitude, 64)
		gisdata.Latitude, _ = strconv.ParseFloat(tmpgisdata[index].Latitude, 64)
		gislist = append(gislist, gisdata)
	}

	h.Data["City"] = cityname
	h.Data["Gis"] = gislist
	h.TplName = "bdmaptest.tpl"
}

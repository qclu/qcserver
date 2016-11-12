package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHospitalCtl"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"PUT"},
			Params: nil})

}

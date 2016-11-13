package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcAdministratorCtl"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"PUT"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDepartmentCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDepartmentCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDepartmentCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDepartmentCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcDevModelCtl"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"PUT"},
			Params: nil})

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

	beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcHwVersionCtl"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"PUT"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcMethodologyCtl"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"PUT"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"] = append(beego.GlobalControllerRouter["qcserver/controllers:QcSwVersionCtl"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/`,
			AllowHTTPMethods: []string{"PUT"},
			Params: nil})

}

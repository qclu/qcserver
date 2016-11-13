// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"qcserver/controllers"

	"github.com/astaxie/beego"
)

func init() {
	hns := beego.NewNamespace("/hospital",
		beego.NSInclude(
			&controllers.QcHospitalCtl{},
		),
	)
	beego.AddNamespace(hns)
	dns := beego.NewNamespace("/department",
		beego.NSInclude(
			&controllers.QcDepartmentCtl{},
		),
	)
	beego.AddNamespace(dns)
	adminns := beego.NewNamespace("/administrator",
		beego.NSInclude(
			&controllers.QcAdministratorCtl{},
		),
	)
	beego.AddNamespace(adminns)
	mthns := beego.NewNamespace("/methodology",
		beego.NSInclude(
			&controllers.QcMethodologyCtl{},
		),
	)
	beego.AddNamespace(mthns)
	devmodelns := beego.NewNamespace("/devmodel",
		beego.NSInclude(
			&controllers.QcDevModelCtl{},
		),
	)
	beego.AddNamespace(devmodelns)
	hwvns := beego.NewNamespace("/hwversion",
		beego.NSInclude(
			&controllers.QcHwVersionCtl{},
		),
	)
	beego.AddNamespace(hwvns)
	swvns := beego.NewNamespace("/swversion",
		beego.NSInclude(
			&controllers.QcSwVersionCtl{},
		),
	)
	beego.AddNamespace(swvns)
	devrelns := beego.NewNamespace("/devrel",
		beego.NSInclude(
			&controllers.QcDevRelCtl{},
		),
	)
	beego.AddNamespace(devrelns)
	regmodelns := beego.NewNamespace("/reagentmodel",
		beego.NSInclude(
			&controllers.QcReagentModelCtl{},
		),
	)
	beego.AddNamespace(regmodelns)
}

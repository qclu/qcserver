package main

import (
	//"btest/controllers"
	//_ "btest/routers"
	//"common/beego"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"qcserver/models"
	_ "qcserver/routers"
	"qcserver/tcpnetwork"
	"qcserver/util/log"
	"time"
)

func main_start() {
	fmt.Println("start...")
	orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8", 30)
	orm.RegisterModel(new(models.QcAdministrator))
	orm.RunSyncdb("default", false, true)

	fmt.Println("Test Administrator...")
	o := orm.NewOrm()
	o.Using("default")

	admin := models.QcAdministrator{
		Username: "qcl",
		Role:     0,
		Password: "1234qwerasasdd",
	}

	admin.Id = 1
	num, err := o.Update(&admin, "Password", "Updated")
	fmt.Println("Num: ", num, " err: ", err)
}

func main_admin() {
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	for i := 0; i < 10; i++ {
		username := fmt.Sprintf("ad_%v", i)
		admin := models.QcAdministrator{
			Username: username,
			Role:     i % 2,
			Password: "1234",
		}

		err = dbSync.InsertQcAdmin(&admin)
		if err != nil {
			logger.LogError("Failed to insert new admin, error: ", err)
		}
	}

	logger.LogInfo("Get all admins info...")
	admins, err := dbSync.GetQcAdmins(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all admins, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(admins); i++ {
		fmt.Println("admin info: ", admins[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Update admin passwd...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(admins); i++ {
		new_passwd := fmt.Sprintf("%s_new", admins[i].Password)
		err = admins[i].UpdatePasswd(dbSync, new_passwd)
		if err != nil {
			logger.LogError("Failed to update passwd, error: ", err)
		}
		new_ad, err := dbSync.GetQcAdmin(admins[i].Username)
		if err != nil {
			logger.LogError("Failed to get admin info, error: ", err)
		}
		fmt.Println("admin: [", new_ad, "]")
	}
	fmt.Println("-------------------------------------------------------------------------")

	admin_0, err := dbSync.GetQcAdmins(0, 0, "")
	if err != nil {
		logger.LogError("Failed to get admins of role 0, error: ", err)
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(admin_0); i++ {
		fmt.Println("admin info: ", admin_0[i])
		err = admin_0[i].Delete(dbSync)
		if err != nil {
			logger.LogError("Failed to delete admin[", admin_0[i], "], error: ", err)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")
}

func main_hospital1() {
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("hospital_%v", i)
		addr := fmt.Sprintf("beijing_%v", i+10)
		gis := fmt.Sprintf("76.23145%v, 89.90897%v3", i%2, i%5)
		h := models.QcHospital{
			Name:    name,
			Addr:    addr,
			Gis:     gis,
			Created: time.Now().Format("2006-01-02 15:04:05"),
			Updated: time.Now().Format("2006-01-02 15:04:05"),
		}

		err = dbSync.InsertQcHospital(&h)
		if err != nil {
			logger.LogError("Failed to insert new hospital, error: ", err)
		}
	}

	logger.LogInfo("Get all hospital info...")
	hospitals, err := dbSync.GetQcHospitals(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hospitals); i++ {
		fmt.Println("hospital info: ", hospitals[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	logger.LogInfo("Get hospital pages(pgidx: 0, pgsize: 3)...")
	hospitals, err = dbSync.GetQcHospitals(0, 3, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hospitals); i++ {
		fmt.Println("hospital info: ", hospitals[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	logger.LogInfo("Get hospital pages(pgidx: 2, pgsize: 3)...")
	hospitals, err = dbSync.GetQcHospitals(2, 3, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hospitals); i++ {
		fmt.Println("hospital info: ", hospitals[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Update hospital addr...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hospitals); i++ {
		new_addr := fmt.Sprintf("%s_new", hospitals[i].Addr)
		err = hospitals[i].UpdateAddr(dbSync, new_addr)
		if err != nil {
			logger.LogError("Failed to update addr, error: ", err)
		}
		new_h, err := dbSync.GetQcHospital(hospitals[i].Name)
		if err != nil {
			logger.LogError("Failed to get hospital info, error: ", err)
		}
		fmt.Println("hospital: [", new_h, "]")
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hospitals); i++ {
		new_addr := fmt.Sprintf("%s_new", hospitals[i].Addr)
		err = hospitals[i].UpdateAddr(dbSync, new_addr)
		if err != nil {
			logger.LogError("Failed to update addr, error: ", err)
		}
		new_h, err := dbSync.GetQcHospital(hospitals[i].Name)
		if err != nil {
			logger.LogError("Failed to get hospital info, error: ", err)
		}
		fmt.Println("hospital: [", new_h, "]")
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("delete hospitals...")
	for i := 0; i < len(hospitals); i++ {
		if i%2 == 0 {
			err = hospitals[i].Delete(dbSync)
			if err != nil {
				logger.LogError("Failed to delete hospital, error: ", err)
			}
			continue
		}
		err := dbSync.DeleteQcHospital(hospitals[i])
		if err != nil {
			logger.LogError("Failed to delete hospital info, error: ", err)
		}
	}
}

func main_devmodel() {
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	fmt.Println("Create methodology...")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("methodology_%v", i)
		anno := fmt.Sprintf("annotation_%v", i*10)
		mt := models.QcMethodology{
			Name:       name,
			Annotation: anno,
			Created:    time.Now().Format(models.TIME_FMT),
			Updated:    time.Now().Format(models.TIME_FMT),
		}

		err = dbSync.InsertQcMethodology(&mt)
		if err != nil {
			logger.LogError("Failed to insert new methodology, error: ", err)
		}
	}

	logger.LogInfo("Get all methodology info...")
	mths, err := dbSync.GetQcMethodologys(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all Methodology, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(mths); i++ {
		fmt.Println("methodology info: ", mths[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Update mths anno ...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(mths); i++ {
		new_anno := fmt.Sprintf("%s_new", mths[i].Annotation)
		err = mths[i].UpdateAnnotation(dbSync, new_anno)
		if err != nil {
			logger.LogError("Failed to update annotation, error: ", err)
		}
		new_mth, err := dbSync.GetQcMethodology(mths[i].Name)
		if err != nil {
			logger.LogError("Failed to get methodology info, error: ", err)
		}
		fmt.Println("methodology: [", new_mth, "]")
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Create Device Model...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("devmodel_%v", i)
		model := fmt.Sprintf("0XNBC%vOCXX%v", i*3, i)
		release := time.Now().Format(models.TIME_FMT)
		anno := fmt.Sprintf("%v_anno", i)
		devmodel, err := models.CreateQcDevModel(dbSync, name, model, release, mths[i], anno)
		if err != nil {
			logger.LogError("Failed to create device model, error: ", err)
		} else {
			fmt.Println("devmodel: ", devmodel)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

	devmodels, err := dbSync.GetQcDevmodels(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list device model info, err: ", err)
	}
	for i := 0; i < len(devmodels); i++ {
		fmt.Println("device model[name: ", devmodels[i].Name, ", Model: ", devmodels[i].Model, ", Methodology: ", devmodels[i].Methodology)
	}

	for i := 0; i < len(mths); i++ {
		mths[i].Delete(dbSync)
	}
}

func main_devmodel1() {
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	fmt.Println("Create methodology...")
	name := fmt.Sprintf("methodology_%v", 101)
	anno := fmt.Sprintf("annotation_%v", 10)
	mt := models.QcMethodology{
		Name:       name,
		Annotation: anno,
		Created:    time.Now().Format(models.TIME_FMT),
		Updated:    time.Now().Format(models.TIME_FMT),
	}

	err = dbSync.InsertQcMethodology(&mt)
	if err != nil {
		logger.LogError("Failed to insert new methodology, error: ", err)
	}

	logger.LogInfo("Get all methodology info...")
	mths, err := dbSync.GetQcMethodologys(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all Methodology, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(mths); i++ {
		fmt.Println("methodology info: ", mths[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Create Device Model...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("devmodel_%v", i)
		model := fmt.Sprintf("0XNBC%vOCXX%v", i*3, i)
		release := time.Now().Format(models.TIME_FMT)
		anno := fmt.Sprintf("%v_anno", i)
		devmodel, err := models.CreateQcDevModel(dbSync, name, model, release, mths[0], anno)
		if err != nil {
			logger.LogError("Failed to create device model, error: ", err)
		} else {
			fmt.Println("devmodel: ", devmodel)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

	logger.LogInfo("Get dev model with mname...")
	dmodel, err := dbSync.GetQcMethodology(mths[0].Name)
	if err != nil {
		logger.LogError("Failed to get dev model cnt of methodology, err: ", err)
	} else {
		logger.LogInfo("dev models[", dmodel, "] belongs to methodology[", mths[0], "]")
	}

	logger.LogInfo("Get dev model of methodology with mname...")
	cnt, err := dbSync.GetDevmodelCntOfMethodologyName(mths[0].Name)
	if err != nil {
		logger.LogError("Failed to get dev model cnt of methodology, err: ", err)
	} else {
		logger.LogInfo("Total ", cnt, "dev models belongs to methodology[", mths[0], "]")
	}

	logger.LogInfo("Get dev model of methodology with mid...")
	cnt, err = dbSync.GetDevmodelCntOfMethodologyId(mths[0].Id)
	if err != nil {
		logger.LogError("Failed to get dev model cnt of methodology, err: ", err)
	} else {
		logger.LogInfo("Total ", cnt, "dev models belongs to methodology[", mths[0], "]")
	}

	devmodels, err := dbSync.GetQcDevmodels(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list device model info, err: ", err)
	}
	for i := 0; i < len(devmodels); i++ {
		fmt.Println("device model[name: ", devmodels[i].Name, ", Model: ", devmodels[i].Model, ", Methodology: ", devmodels[i].Methodology)
	}

	//for i := 0; i < len(mths); i++ {
	//	mths[i].Delete(dbSync)
	//}
}

func main_shwversion() {
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	fmt.Println("Create methodology...")
	name := fmt.Sprintf("methodology_%v", 101)
	anno := fmt.Sprintf("annotation_%v", 10)
	mt := models.QcMethodology{
		Name:       name,
		Annotation: anno,
		Created:    time.Now().Format(models.TIME_FMT),
		Updated:    time.Now().Format(models.TIME_FMT),
	}

	err = dbSync.InsertQcMethodology(&mt)
	if err != nil {
		logger.LogError("Failed to insert new methodology, error: ", err)
	}

	logger.LogInfo("Get all methodology info...")
	mths, err := dbSync.GetQcMethodologys(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all Methodology, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(mths); i++ {
		fmt.Println("methodology info: ", mths[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Create Device Model...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("devmodel_%v", i)
		model := fmt.Sprintf("0XNBC%vOCXX%v", i*3, i)
		release := time.Now().Format(models.TIME_FMT)
		anno := fmt.Sprintf("%v_anno", i)
		devmodel, err := models.CreateQcDevModel(dbSync, name, model, release, mths[0], anno)
		if err != nil {
			logger.LogError("Failed to create device model, error: ", err)
		} else {
			fmt.Println("devmodel: ", devmodel)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("Get all dev models...")
	devmodels, err := dbSync.GetQcDevmodels(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list device model info, err: ", err)
	}
	for i := 0; i < len(devmodels); i++ {
		fmt.Println("device model[name: ", devmodels[i].Name, ", Model: ", devmodels[i].Model, ", Methodology: ", devmodels[i].Methodology, "]")
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Create hardware version...")
	models_cnt := len(devmodels)
	for i := 0; i < 20; i++ {
		dmidx := rand.Intn(models_cnt)
		devModel := devmodels[dmidx]
		version := fmt.Sprintf("0XNBC%vOCXX%v", i*3, i)
		anno := fmt.Sprintf("%v_anno", i)
		hwversion, err := models.CreateQcHwVersion(dbSync, devModel, version, anno)
		if err != nil {
			logger.LogError("Failed to create hardware version, error: ", err)
		} else {
			fmt.Println("hwversion: ", hwversion)
		}
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Get all hardware version...")
	hwversions, err := dbSync.GetQcHwVersions(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list hardware version info, err: ", err)
	}
	for i := 0; i < len(hwversions); i++ {
		fmt.Println("hardware version[DevModel: ", hwversions[i].DevModel,
			", Version: ", hwversions[i].Version,
			", Anno: ", hwversions[i].Annotation,
			", Created: ", hwversions[i].Created,
			", Updated: ", hwversions[i].Updated, "]")
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Create software version...")
	hwversion_cnt := len(hwversions)
	for i := 0; i < 20; i++ {
		smidx := rand.Intn(hwversion_cnt)
		hwv := hwversions[smidx]
		version := fmt.Sprintf("0XNBC%vOCXX%v", i*3, i)
		swtype := "RELEASE"
		desp := fmt.Sprintf("%v_desp", i)
		swversion, err := models.CreateQcSwVersion(dbSync, hwv, version, swtype, desp)
		if err != nil {
			logger.LogError("Failed to create software version, error: ", err)
		} else {
			fmt.Println("swversion: ", swversion)
		}
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Get all software version...")
	swversions, err := dbSync.GetQcSwVersions(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list software version info, err: ", err)
	}
	for i := 0; i < len(swversions); i++ {
		fmt.Println("software version[HwVersion: ", swversions[i].HwVersion,
			", Version: ", swversions[i].Version,
			", SwType: ", swversions[i].SwType,
			", Description: ", swversions[i].Description,
			", Created: ", hwversions[i].Created,
			", Updated: ", hwversions[i].Updated, "]")
	}
}

func main_hospital() {
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	fmt.Println("Create hospital...")
	name := fmt.Sprintf("hospital_%v", 101)
	addr := fmt.Sprintf("Beijing_%v", 10)
	gis := fmt.Sprintf("%v, %v", 123.12, 12.123)
	h := models.QcHospital{
		Name:    name,
		Addr:    addr,
		Gis:     gis,
		Created: time.Now().Format(models.TIME_FMT),
		Updated: time.Now().Format(models.TIME_FMT),
	}

	err = dbSync.InsertQcHospital(&h)
	if err != nil {
		logger.LogError("Failed to insert new hospital, error: ", err)
	}

	logger.LogInfo("Get all methodology info...")
	hs, err := dbSync.GetQcHospitals(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hs); i++ {
		fmt.Println("hospital info: ", hs[i])
	}
	fmt.Println("-------------------------------------------------------------------------")
}

func main_department() {
	fmt.Println("start ...")
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	fmt.Println("Create hospital...")
	name := fmt.Sprintf("hospital_%v", 101)
	addr := fmt.Sprintf("Beijing_%v", 10)
	gis := fmt.Sprintf("%v, %v", 123.12, 12.123)
	h := models.QcHospital{
		Name:    name,
		Addr:    addr,
		Gis:     gis,
		Created: time.Now().Format(models.TIME_FMT),
		Updated: time.Now().Format(models.TIME_FMT),
	}

	err = dbSync.InsertQcHospital(&h)
	if err != nil {
		logger.LogError("Failed to insert new hospital, error: ", err)
	}

	logger.LogInfo("Get all hospital info...")
	hs, err := dbSync.GetQcHospitals(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hs); i++ {
		fmt.Println("hospital info: ", hs[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("Create department...")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("department_%v", i)
		department, err := models.CreateQcDepartment(dbSync, name, hs[i%len(hs)])
		if err != nil {
			logger.LogError("Failed to create department, error: ", err)
		} else {
			fmt.Println("department: ", department)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("Get department info...")
	department, err := dbSync.GetQcDepartment("department_1", hs[0].Name)
	if err != nil {
		logger.LogError("Failed to get department[department_1] of hospital[", hs[0].Name, "], err: ", err)
	} else {
		fmt.Println("Department info: ", department, "   ", department.Hospital.Name)
	}
	fmt.Println("Delete department ...")
	//err = dbSync.DeleteQcDepartmentSQL(hs[0]., "department_1")
	fmt.Println("err: ", err)
}

func main_devrel() {
	fmt.Println("start ...")
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")
	dbSync, err := models.NewDBSync("mysql", "root:root@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to init database module, error:", err)
		return
	}

	fmt.Println("Create hospital...")
	name := fmt.Sprintf("hospital_%v", 101)
	addr := fmt.Sprintf("Beijing_%v", 10)
	gis := fmt.Sprintf("%v, %v", 123.12, 12.123)
	h := models.QcHospital{
		Name:    name,
		Addr:    addr,
		Gis:     gis,
		Created: time.Now().Format(models.TIME_FMT),
		Updated: time.Now().Format(models.TIME_FMT),
	}

	err = dbSync.InsertQcHospital(&h)
	if err != nil {
		logger.LogError("Failed to insert new hospital, error: ", err)
	}

	logger.LogInfo("Get all hospital info...")
	hs, err := dbSync.GetQcHospitals(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(hs); i++ {
		fmt.Println("hospital info: ", hs[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("Create department...")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("department_%v", i)
		department, err := models.CreateQcDepartment(dbSync, name, hs[i%len(hs)])
		if err != nil {
			logger.LogError("Failed to create department, error: ", err)
		} else {
			fmt.Println("department: ", department)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("Get department info...")
	department, err := dbSync.GetQcDepartment("department_1", hs[0].Name)
	if err != nil {
		logger.LogError("Failed to get department[department_1] of hospital[", hs[0].Name, "], err: ", err)
	} else {
		fmt.Println("Department info: ", department, "   ", department.Hospital.Name)
	}
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Get all departments...")
	departments, err := dbSync.GetQcDepartments(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all hospitals, error: ", err)
		return
	}
	for i := 0; i < len(departments); i++ {
		fmt.Println("department info: ", departments[i])
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Create methodology...")
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("methodology_%v%v", rand.Intn(10), rand.Intn(10))
		anno := fmt.Sprintf("annotation_%v", rand.Intn(10))
		mt := models.QcMethodology{
			Name:       name,
			Annotation: anno,
			Created:    time.Now().Format(models.TIME_FMT),
			Updated:    time.Now().Format(models.TIME_FMT),
		}
		err = dbSync.InsertQcMethodology(&mt)
		if err != nil {
			logger.LogError("Failed to insert new methodology, error: ", err)
		}
	}

	logger.LogInfo("Get all methodology info...")
	mths, err := dbSync.GetQcMethodologys(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list all Methodology, error: ", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < len(mths); i++ {
		fmt.Println("methodology info: ", mths[i])
	}
	fmt.Println("-------------------------------------------------------------------------")

	time.Sleep(10 * time.Second)

	fmt.Println("Create Device Model...")
	fmt.Println("-------------------------------------------------------------------------")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("devmodel_%v", rand.Intn(10))
		model := fmt.Sprintf("0XNBC%vOCXX%v", rand.Intn(10), rand.Intn(10))
		release := time.Now().Format(models.TIME_FMT)
		anno := fmt.Sprintf("%v_anno", rand.Intn(10))
		devmodel, err := models.CreateQcDevModel(dbSync, name, model, release, mths[i%len(mths)], anno)
		if err != nil {
			logger.LogError("Failed to create device model, error: ", err)
		} else {
			fmt.Println("devmodel: ", devmodel)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

	fmt.Println("Get all dev models...")
	devmodels, err := dbSync.GetQcDevmodels(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list device model info, err: ", err)
	}
	for i := 0; i < len(devmodels); i++ {
		fmt.Println("device model[name: ", devmodels[i].Name, ", Model: ", devmodels[i].Model, ", Methodology: ", devmodels[i].Methodology, "]")
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Create hardware version...")
	models_cnt := len(devmodels)
	for i := 0; i < 20; i++ {
		dmidx := rand.Intn(models_cnt)
		devModel := devmodels[dmidx]
		version := fmt.Sprintf("0XNBC%vOCXX%v", rand.Intn(10), rand.Intn(10))
		anno := fmt.Sprintf("%v_anno", rand.Intn(10))
		hwversion, err := models.CreateQcHwVersion(dbSync, devModel, version, anno)
		if err != nil {
			logger.LogError("Failed to create hardware version, error: ", err)
		} else {
			fmt.Println("hwversion: ", hwversion)
		}
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Get all hardware version...")
	hwversions, err := dbSync.GetQcHwVersions(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list hardware version info, err: ", err)
	}
	for i := 0; i < len(hwversions); i++ {
		fmt.Println("hardware version[DevModel: ", hwversions[i].DevModel,
			", Version: ", hwversions[i].Version,
			", Anno: ", hwversions[i].Annotation,
			", Created: ", hwversions[i].Created,
			", Updated: ", hwversions[i].Updated, "]")
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Create software version...")
	hwversion_cnt := len(hwversions)
	for i := 0; i < 20; i++ {
		smidx := rand.Intn(hwversion_cnt)
		hwv := hwversions[smidx]
		version := fmt.Sprintf("0XNBC%vOCXX%v", rand.Intn(10), rand.Intn(10))
		swtype := "RELEASE"
		desp := fmt.Sprintf("%v_desp", rand.Intn(10))
		swversion, err := models.CreateQcSwVersion(dbSync, hwv, version, swtype, desp)
		if err != nil {
			logger.LogError("Failed to create software version, error: ", err)
		} else {
			fmt.Println("swversion: ", swversion)
		}
	}

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Get all software version...")
	swversions, err := dbSync.GetQcSwVersions(0, 0, "")
	if err != nil {
		logger.LogError("Failed to list software version info, err: ", err)
	}
	for i := 0; i < len(swversions); i++ {
		fmt.Println("software version[HwVersion: ", swversions[i].HwVersion,
			", Version: ", swversions[i].Version,
			", SwType: ", swversions[i].SwType,
			", Description: ", swversions[i].Description,
			", Created: ", hwversions[i].Created,
			", Updated: ", hwversions[i].Updated, "]")
	}
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Get all software version...")

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Create device release info...")
	for i := 0; i < 10; i++ {
		sn := fmt.Sprintf("0xBC%vDXE%vAOO%v", i+3, i+1, i)
		//smcard := fmt.Sprintf("SMCARD_%v", i)
		department, err := models.CreateQcDevRel(dbSync, sn,
			time.Now().Format(models.TIME_FMT), swversions[i%len(swversions)], departments[i%len(departments)])
		if err != nil {
			logger.LogError("Failed to create department, error: ", err)
		} else {
			fmt.Println("department: ", department)
		}
	}
	fmt.Println("-------------------------------------------------------------------------")

}

func main() {
	fmt.Println("start ...")
	defer time.Sleep(time.Second)
	logger, err := log.NewLog("/var/log/", "qcserver", 0)
	if err != nil {
		fmt.Println("Failed to init log module...")
		return
	}
	logger.LogInfo("Info: log module start...")

	//err = models.DBSyncInit("mysql", "root:1qaz@WSX@/qcdatabase?charset=utf8")
	err = models.DBSyncInit("mysql", "root:123qwe@/orm_test?charset=utf8")
	if err != nil {
		logger.LogError("Failed to start database process module, err:", err)
		return
	}

	//start tcp service
	//tcpserver := tcpnetwork.NewTcpServer("118.178.188.139:14444")
	tcpserver := tcpnetwork.NewTcpServer("127.0.0.1:14444")
	err = tcpserver.Listen()
	if err != nil {
		logger.LogError("Failed to start tcpserver listening routine, err:", err)
		return
	}
	go tcpserver.Start()

	beego.BConfig.RunMode = "dev"
	beego.SetViewsPath("views")
	beego.SetStaticPath("/swpackage", "/data")
	beego.Run()
}

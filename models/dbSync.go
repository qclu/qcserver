package models

import (
	//"github.com/beego/orm"
	"github.com/astaxie/beego/orm"
	//"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	"qcserver/util/log"
	con "strconv"
	"sync"
	"time"
)

type DBSync struct {
	logger *log.Log
	mutex  sync.Mutex
}

const RetryTime = 3
const RetrySleepDuration = 10 * time.Millisecond

var dbSync *DBSync

func DBSyncInit(dbDriver, dbDataSource string) error {
	var err error
	dbSync, err = NewDBSync(dbDriver, dbDataSource)
	if err != nil {
		dbSync = nil
		fmt.Println("Failed to init database driver, err: ", err)
	}
	return err
}

func GetDBSync() *DBSync {
	return dbSync
}

func NewDBSync(dbDriver, dbDataSource string) (*DBSync, error) {
	orm.RegisterDriver(dbDriver, orm.DRMySQL)
	database := "default"
	//fmt.Println(dbDriver,dbDataSource)
	orm.RegisterDataBase(database, dbDriver, dbDataSource)

	orm.RegisterModel(new(QcAdministrator))
	orm.RegisterModel(new(QcHospital))
	orm.RegisterModel(new(QcMethodology))
	orm.RegisterModel(new(QcDevModel))
	orm.RegisterModel(new(QcHwVersion))
	orm.RegisterModel(new(QcSwVersion))
	orm.RegisterModel(new(QcDepartment))
	orm.RegisterModel(new(QcDevRel))
	//orm.RegisterModel(new(NodesGroup))

	orm.SetMaxOpenConns(database, 10)
	orm.SetMaxIdleConns(database, 10)
	//create table
	forceCreate := false
	verbose := true
	err := orm.RunSyncdb(database, forceCreate, verbose)
	if err != nil {
		return nil, err
	}

	logger := log.GetLog()
	println("log:", logger)
	return &DBSync{logger: logger}, nil
}

func (d *DBSync) GetPagesInfo(tableName string, currentpage int, pagesize int, conditions string) (int, int, orm.RawSeter) {
	if currentpage <= 0 {
		currentpage = 0
	}
	var rs orm.RawSeter
	o := orm.NewOrm()
	var totalItem, totalpages int = 0, 0                                                          //×ÜÌõÊý,×ÜÒ³Êý
	o.Raw("SELECT count(*) FROM " + tableName + "  where 1>0 " + conditions).QueryRow(&totalItem) //»ñÈ¡×ÜÌõÊý
	if pagesize == 0 {
		pagesize = totalItem
	}
	if totalItem <= pagesize {
		totalpages = 1
	} else if totalItem > pagesize {
		temp := totalItem / pagesize
		if (totalItem % pagesize) != 0 {
			temp = temp + 1
		}
		totalpages = temp
	}
	rs = o.Raw("select *  from  " + tableName + "  where id >0 " + conditions + " LIMIT " + con.Itoa((currentpage)*pagesize) + "," + con.Itoa(pagesize))
	return totalItem, totalpages, rs
}

func (d *DBSync) InsertQcAdmin(admin *QcAdministrator) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(admin)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcDevModel(obj *QcDevModel) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcDevRel(obj *QcDevRel) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcSwVersion(obj *QcSwVersion) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcHwVersion(obj *QcHwVersion) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcMethodology(m *QcMethodology) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(m)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcDepartment(department *QcDepartment) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(department)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) InsertQcHospital(hospital *QcHospital) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(hospital)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcAdmin(admin *QcAdministrator) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(admin)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcSwVersion(obj *QcSwVersion) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcHwVersion(obj *QcHwVersion) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcDevRel(obj *QcDevRel) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcDevModel(obj *QcDevModel) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcMethdology(m *QcMethodology) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(m)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcDepartment(department *QcDepartment) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(department)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) DeleteQcHospital(hospital *QcHospital) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(hospital)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) GetQcAdmin(username string) (*QcAdministrator, error) {
	params := map[string]interface{}{"Username": username}
	var admin QcAdministrator
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_ADMINISTRATOR)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&admin)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (d *DBSync) GetQcSwVersion(version string) (*QcSwVersion, error) {
	params := map[string]interface{}{"Version": version}
	var obj QcSwVersion
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_SWVERSION)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&obj)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *DBSync) GetQcHwVersion(version string) (*QcHwVersion, error) {
	params := map[string]interface{}{"Version": version}
	var obj QcHwVersion
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_HWVERSION)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&obj)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *DBSync) GetQcDevModel(name string) (*QcDevModel, error) {
	params := map[string]interface{}{"Name": name}
	var obj QcDevModel
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_DEVMODEL)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&obj)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *DBSync) GetQcDevRel(sn string) (*QcDevRel, error) {
	params := map[string]interface{}{"Sn": sn}
	var obj QcDevRel
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_DEVREL)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&obj)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *DBSync) GetQcMethodology(name string) (*QcMethodology, error) {
	params := map[string]interface{}{"Name": name}
	var m QcMethodology
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_METHODOLOGY)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&m)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *DBSync) GetQcDepartmentsWithHospital(hname string) ([]*QcDepartment, error) {
	var objs []*QcDepartment
	var err error
	hospital, err := d.GetQcHospital(hname)
	if err != nil {
		return nil, err
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	_, err = ormer.QueryTable(DB_T_DEPARTMENT).Filter("Hospital", hospital.Id).RelatedSel().All(&objs)
	if err != nil {
		d.logger.LogError("Failed to list all departments of hospital[", hospital, "], error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcDepartments() ([]*QcDepartment, error) {
	var objs []*QcDepartment
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_DEPARTMENT)
	if _, err = qs.All(&objs); err != nil {
		d.logger.LogError("Failed to list all departments, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcHospital(name string) (*QcHospital, error) {
	params := map[string]interface{}{"Name": name}
	var hospital QcHospital
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_HOSPITAL)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&hospital)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (d *DBSync) GetQcDepartment(dname, hname string) (*QcDepartment, error) {
	var department QcDepartment
	var err error
	hospital, err := d.GetQcHospital(hname)
	if err != nil {
		return nil, err
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	for retry := 0; retry < RetryTime; retry++ {
		err = ormer.QueryTable(DB_T_DEPARTMENT).Filter("Name", dname).Filter("Hospital", hospital.Id).RelatedSel().Limit(1).One(&department)
		if err != nil {
			if err == orm.ErrNoRows {
				return nil, errors.New(ERR_OBJ_NOT_EXIST)
			} else {
				d.logger.LogWarn(err)
				continue
			}
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (d *DBSync) UpdateQcAdmin(admin *QcAdministrator) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	admin.Updated = time.Now().Format(TIME_FMT)
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(admin)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcSwVersion(obj *QcSwVersion) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	obj.Updated = time.Now().Format(TIME_FMT)
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcHwVersion(obj *QcHwVersion) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	obj.Updated = time.Now().Format(TIME_FMT)
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcDevRel(obj *QcDevRel) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	obj.Updated = time.Now().Format(TIME_FMT)
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcDevModel(obj *QcDevModel) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	obj.Updated = time.Now().Format(TIME_FMT)
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(obj)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcMethodology(m *QcMethodology) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	m.Updated = time.Now().Format("2006-01-02 15:04:05")
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(m)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcDepartment(h *QcDepartment) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	h.Updated = time.Now().Format("2006-01-02 15:04:05")
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(h)
		if err == nil {
			return err
		}
	}
	return err
}

func (d *DBSync) UpdateQcHospital(h *QcHospital) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	h.Updated = time.Now().Format("2006-01-02 15:04:05")
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Update(h)
		if err == nil {
			return err
		}
	}
	return err
}

//if role is set to -1, all admins info will be listed
func (d *DBSync) GetQcAdmins(role int) ([]*QcAdministrator, error) {
	var admins []*QcAdministrator
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_ADMINISTRATOR)
	if role != -1 {
		qs = qs.Filter("Role", role)
	}
	if _, err = qs.All(&admins); err != nil {
		d.logger.LogError("Failed to list admins, error: ", err)
		return nil, err
	}
	return admins, nil
}

func (d *DBSync) GetQcSwVersions() ([]*QcSwVersion, error) {
	var objs []*QcSwVersion
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_SWVERSION)
	if _, err = qs.All(&objs); err != nil {
		d.logger.LogError("Failed to list admins, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcHwVersions() ([]*QcHwVersion, error) {
	var objs []*QcHwVersion
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_HWVERSION)
	if _, err = qs.All(&objs); err != nil {
		d.logger.LogError("Failed to list admins, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcDevmodels() ([]*QcDevModel, error) {
	var objs []*QcDevModel
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_DEVMODEL)
	if _, err = qs.All(&objs); err != nil {
		d.logger.LogError("Failed to list admins, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcMethodologys() ([]*QcMethodology, error) {
	var ms []*QcMethodology
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_METHODOLOGY)
	if _, err = qs.All(&ms); err != nil {
		d.logger.LogError("Failed to list methodologies, error: ", err)
		return nil, err
	}
	return ms, nil
}

func (d *DBSync) GetQcHospitals(pgidx, pgsize int, conditions string) ([]*QcHospital, error) {
	var hospitals []*QcHospital
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_HOSPITAL, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&hospitals); err != nil {
		d.logger.LogError("Failed to list hospitals, error: ", err)
		return nil, err
	}
	return hospitals, nil
}

//func (d *DBSync) GetQcHospitals(pgcnt, pgnum int32) ([]*QcHospital, error) {
//	var hospitals []*QcHospital
//	var err error
//	d.mutex.Lock()
//	defer d.mutex.Unlock()
//	ormer := orm.NewOrm()
//	qs := ormer.QueryTable(DB_T_HOSPITAL)
//	if _, err = qs.All(&hospitals); err != nil {
//		d.logger.LogError("Failed to list hospitals, error: ", err)
//		return nil, err
//	}
//	return hospitals, nil
//}

func (d *DBSync) GetDevmodelCntOfMethodologyId(mid int64) (int64, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	qs := ormer.QueryTable(DB_T_DEVMODEL)
	return qs.Filter("methodology_id", mid).Count()
}

func (d *DBSync) GetDevmodelCntOfMethodologyName(mth_name string) (int64, error) {
	mth, err := d.GetQcMethodology(mth_name)
	if err != nil {
		return 0, err
	}
	return d.GetDevmodelCntOfMethodologyId(mth.Id)
}

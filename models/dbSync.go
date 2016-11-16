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
	orm.RegisterModel(new(QcReagentModel))
	orm.RegisterModel(new(QcReagentProduce))
	orm.RegisterModel(new(QcReagentRel))
	orm.RegisterModel(new(QcQcProduct))

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

func (d *DBSync) InsertQcQcLotnum(obj *QcQcLotnum) error {
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

func (d *DBSync) InsertQcQcProduct(obj *QcQcProduct) error {
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

func (d *DBSync) InsertQcReagentRel(obj *QcReagentRel) error {
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

func (d *DBSync) InsertQcReagentProduce(obj *QcReagentProduce) error {
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

func (d *DBSync) InsertQcReagentModel(obj *QcReagentModel) error {
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

func (d *DBSync) InsertQcDevLog(log_ent *QcDevLog) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Insert(log_ent)
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

func (d *DBSync) DeleteQcHospitalSQL(name string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_HOSPITAL + " where name='" + name + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcMethdologySQL(name string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_METHODOLOGY + " where name='" + name + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcAdminSQL(username string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_ADMINISTRATOR + " where username='" + username + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
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

func (d *DBSync) DeleteQcSwVersionSQL(version string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_SWVERSION + " where version='" + version + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
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

func (d *DBSync) DeleteQcQcLotnum(obj *QcQcLotnum) error {
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

func (d *DBSync) DeleteQcQcProduct(obj *QcQcProduct) error {
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

func (d *DBSync) DeleteQcReagentRel(obj *QcReagentRel) error {
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

func (d *DBSync) DeleteQcReagentProduce(obj *QcReagentProduce) error {
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

func (d *DBSync) DeleteQcReagentModel(obj *QcReagentModel) error {
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

func (d *DBSync) DeleteQcReagentRelSQL(serial string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_REGREL + " where release_serial='" + serial + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcQcProductSQL(name string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_QCPRODUCT + " where name='" + name + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcReagentProduceSQL(serial string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_REGPRODUCE + " where serial_num='" + serial + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcReagentModelSQL(name string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_REGMODEL + " where name='" + name + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcDevRelSQL(sn string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_DEVREL + " where sn='" + sn + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcHwVersionSQL(version string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_HWVERSION + " where version='" + version + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	return err
}

func (d *DBSync) DeleteQcDevModelSQL(name string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_DEVMODEL + " where name='" + name + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
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

func (d *DBSync) DeleteQcDepartmentSQL(hname, dname string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	sql := "delete from " + DB_T_DEPARTMENT + " using " + DB_T_DEPARTMENT + ", " + DB_T_HOSPITAL +
		" where " + DB_T_DEPARTMENT + ".hospital_id=" + DB_T_HOSPITAL + ".id and " + DB_T_HOSPITAL + ".name='" + hname + "'"
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	_, err := o.Raw(sql).Exec()
	d.logger.LogInfo("delete department")
	return err
}

func (d *DBSync) DeleteQcDepartment(department *QcDepartment) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	dbSync.logger.LogInfo("delete department: ", department)
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		_, err = ormer.Delete(department)
		if err == nil {
			return err
		}
	}
	dbSync.logger.LogInfo("delete department: ", department)
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

func (d *DBSync) GetQcDevmodel(name string) (*QcDevModel, error) {
	params := map[string]interface{}{"name": name}
	var devmodel QcDevModel
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_DEVMODEL)
		for k, v := range params {
			qs = qs.Filter(k, v)
		}
		err = qs.One(&devmodel)
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
	return &devmodel, nil
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

func (d *DBSync) GetQcQcLotnum(id int64) (*QcQcLotnum, error) {
	params := map[string]interface{}{"Id": id}
	var obj QcQcLotnum
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_QCLOTNUM)
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

func (d *DBSync) GetQcQcProducts(pgidx, pgsize int, conditions string) ([]*QcQcProduct, error) {
	var objs []*QcQcProduct
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_QCPRODUCT, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list qcproducts, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcQcProduct(name string) (*QcQcProduct, error) {
	params := map[string]interface{}{"Name": name}
	var obj QcQcProduct
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_QCPRODUCT)
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

func (d *DBSync) GetQcReagentRel(rel_serial string) (*QcReagentRel, error) {
	params := map[string]interface{}{"ReleaseSerial": rel_serial}
	var obj QcReagentRel
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_REGREL)
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

func (d *DBSync) GetQcReagentProduce(serial string) (*QcReagentProduce, error) {
	params := map[string]interface{}{"SerialNum": serial}
	var obj QcReagentProduce
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_REGPRODUCE)
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

func (d *DBSync) GetQcReagentRels(pgidx, pgsize int, conditions string) ([]*QcReagentProduce, error) {
	var objs []*QcReagentProduce
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_REGREL, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list reagent release, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcReagentProduces(pgidx, pgsize int, conditions string) ([]*QcReagentProduce, error) {
	var objs []*QcReagentProduce
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_REGPRODUCE, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list reagent products, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcReagentModels(pgidx, pgsize int, conditions string) ([]*QcReagentModel, error) {
	var objs []*QcReagentModel
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_REGMODEL, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list reagent models, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcReagentModel(name string) (*QcReagentModel, error) {
	params := map[string]interface{}{"Name": name}
	var obj QcReagentModel
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	var err error
	for retry := 0; retry < RetryTime; retry++ {
		qs := ormer.QueryTable(DB_T_REGMODEL)
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
	d.logger.LogInfo("Get department[hospital: ", hname, " department: ", dname, "]")
	var department QcDepartment
	var err error
	hospital, err := d.GetQcHospital(hname)
	if err != nil {
		return nil, err
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ormer := orm.NewOrm()
	err = ormer.QueryTable(DB_T_DEPARTMENT).Filter("Name", dname).Filter("Hospital", hospital.Id).RelatedSel().Limit(1).One(&department)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.New(ERR_OBJ_NOT_EXIST)
		} else {
			d.logger.LogWarn(err)
		}
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

func (d *DBSync) UpdateQcReagentProduce(obj *QcReagentProduce) error {
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

func (d *DBSync) UpdateQcQcProduct(obj *QcQcProduct) error {
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

func (d *DBSync) UpdateQcReagentRel(obj *QcReagentRel) error {
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

func (d *DBSync) UpdateQcReagentModel(obj *QcReagentModel) error {
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
func (d *DBSync) GetQcAdmins(pgidx, pgsize int, conditions string) ([]*QcAdministrator, error) {
	var admins []*QcAdministrator
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_ADMINISTRATOR, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&admins); err != nil {
		d.logger.LogError("Failed to list admins, error: ", err)
		return nil, err
	}
	return admins, nil
}

func (d *DBSync) GetQcSwVersions(pgidx, pgsize int, conditions string) ([]*QcSwVersion, error) {
	var objs []*QcSwVersion
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_SWVERSION, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list swversions, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcHwVersions(pgidx, pgsize int, conditions string) ([]*QcHwVersion, error) {
	var objs []*QcHwVersion
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_HWVERSION, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list devmodels, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetTotalCnt(table string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	var totalcnt int
	totalcnt = 0
	sql := "select count(*) from " + table
	o := orm.NewOrm()
	d.logger.LogInfo(sql)
	err := o.Raw(sql).QueryRow(&totalcnt)
	if err != nil {
		d.logger.LogError("Failed to get elements count from ", table, ", error: ", err)
		return 0, err
	}
	return totalcnt, nil
}

func (d *DBSync) GetQcDepartments(pgidx, pgsize int, conditions string) ([]*QcDepartment, error) {
	var objs []*QcDepartment
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_DEPARTMENT, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list devmodels, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcDevmodels(pgidx, pgsize int, conditions string) ([]*QcDevModel, error) {
	var objs []*QcDevModel
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_DEVMODEL, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&objs); err != nil {
		d.logger.LogError("Failed to list devmodels, error: ", err)
		return nil, err
	}
	return objs, nil
}

func (d *DBSync) GetQcDevRels(pgidx, pgsize int, conditions string) ([]*QcDevRel, error) {
	var ms []*QcDevRel
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_DEVREL, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&ms); err != nil {
		d.logger.LogError("Failed to list devrels, error: ", err)
		return nil, err
	}
	return ms, nil
}

func (d *DBSync) GetQcMethodologys(pgidx, pgsize int, conditions string) ([]*QcMethodology, error) {
	var ms []*QcMethodology
	var err error
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, _, qs := d.GetPagesInfo(DB_T_METHODOLOGY, pgidx, pgsize, conditions)
	if _, err = qs.QueryRows(&ms); err != nil {
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

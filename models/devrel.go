package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcDevRel struct {
	Id        int64        `orm: "pk;auto"`
	SwVersion *QcSwVersion `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Sn        string       `orm:"size(256);unique"`
	Date      string       `orm:"size(10)"`
	//SmCard    string        `orm:"size(128)"`
	Receiver *QcDepartment `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Created  string        `orm:"size(20)"`
	Updated  string        `orm:"size(20)"`
	mutex    sync.Mutex    `orm:"-"`
}

func CreateQcDevRel(dbSync *DBSync, sn, date string, swv *QcSwVersion, department *QcDepartment) (*QcDevRel, error) {
	obj := &QcDevRel{
		Sn: sn,
		//SmCard:    smcard,
		SwVersion: swv,
		Date:      date,
		Receiver:  department,
		Created:   time.Now().Format(TIME_FMT),
		Updated:   time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcDevRel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new device model to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcDevRel(dbSync *DBSync, sn string) error {
	obj, err := dbSync.GetQcDevRel(sn)
	if err != nil {
		dbSync.logger.LogError("Failed to delete device release[SN: ", sn, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcDevRel) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcDevRel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete device model, error: ", err)
	}
	return err
}

func (obj *QcDevRel) UpdateSn(dbSync *DBSync, sn string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	obj.Sn = sn
	err := dbSync.UpdateQcDevRel(obj)
	return err
}

func (obj *QcDevRel) Update(dbSync *DBSync, swv *QcSwVersion, sn, date, smcard string, recv *QcDepartment) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if swv != nil {
		obj.SwVersion = swv
	}
	if len(sn) == 0 {
		obj.Sn = sn
	}
	if len(date) == 0 {
		obj.Date = date
	}
	//if len(smcard) == 0 {
	//	obj.SmCard = smcard
	//}
	if recv != nil {
		obj.Receiver = recv
	}
	err := dbSync.UpdateQcDevRel(obj)
	return err
}

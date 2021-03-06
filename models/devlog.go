package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcDevLog struct {
	Id       int64      `orm: "pk;auto"`
	Dev      *QcDevRel  `orm:"rel(fk);on_delete(do_nothing)"`
	Type     *QcLogType `orm:"rel(fk);on_delete(do_nothing)"`
	Reported string     `orm:"size(20)"`
	mutex    sync.Mutex `orm:"-"`
}

func CreateQcDevLog(dbSync *DBSync, log_type int, content string, dev *QcDevRel) (*QcDevLog, error) {
	logtype, err := dbSync.GetQcLogTypeWithId(log_type)
	if err != nil {
		dbSync.logger.LogError("Failed to create new dev log as logtype with id ", log_type, " doesnot exist...")
		return nil, err
	}
	obj := &QcDevLog{
		Dev:      dev,
		Type:     logtype,
		Reported: time.Now().Format(TIME_FMT),
	}
	err = dbSync.InsertQcDevLog(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new dev log entry to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

//func DeleteQcDepartment(dbSync *DBSync, dname, hname string) error {
//	department, err := dbSync.GetQcDepartment(dname, hname)
//	if err != nil {
//		dbSync.logger.LogError("Failed to delete department[", dname, "] hospital[", hname, "], error: ", err)
//		return err
//	}
//	err = department.Delete(dbSync)
//	return err
//}
//
//func (h *QcDepartment) Delete(dbSync *DBSync) error {
//	h.mutex.Lock()
//	defer h.mutex.Unlock()
//	err := dbSync.DeleteQcDepartment(h)
//	if err != nil {
//		dbSync.logger.LogError("Failed to delete department, error: ", err)
//	}
//	return err
//}
//
//func (h *QcDepartment) UpdateName(dbSync *DBSync, name string) error {
//	h.mutex.Lock()
//	defer h.mutex.Unlock()
//	h.Name = name
//	err := dbSync.UpdateQcDepartment(h)
//	return err
//}

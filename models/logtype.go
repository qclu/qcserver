package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcLogType struct {
	Id      int64 `orm: "pk;auto"`
	Type    uint16
	Level   uint16
	Updated string     `orm:"size(20)"`
	Content string     `orm:"size(4096)"`
	mutex   sync.Mutex `orm:"-"`
}

func CreateQcLogType(dbSync *DBSync, log_type uint16, lvl uint16, content string) (*QcLogType, error) {
	obj := &QcLogType{
		Type:    log_type,
		Level:   lvl,
		Content: content,
		Updated: time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcLogType(obj)
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

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcDepartment struct {
	Id       int64       `orm: "pk;auto"`
	Name     string      `orm:"size(256);unique"`
	Hospital *QcHospital `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Created  string      `orm:"size(20)"`
	Updated  string      `orm:"size(20)"`
	mutex    sync.Mutex  `orm:"-"`
}

func CreateQcDepartment(dbSync *DBSync, name string, hospital *QcHospital) (*QcDepartment, error) {
	h := &QcDepartment{
		Name:     name,
		Hospital: hospital,
		Created:  time.Now().Format(TIME_FMT),
		Updated:  time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcDepartment(h)
	if err != nil {
		dbSync.logger.LogError("Failed to add new department info to database, error: ", err)
		return nil, err
	}
	return h, nil
}

func DeleteQcDepartment(dbSync *DBSync, dname, hname string) error {
	department, err := dbSync.GetQcDepartment(dname, hname)
	if err != nil {
		dbSync.logger.LogError("Failed to delete hospital[", name, "], error: ", err)
		return err
	}
	err = department.Delete(dbSync)
	return err
}

func (h *QcDepartment) Delete(dbSync *DBSync) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	err := dbSync.DeleteQcDepartment(h)
	if err != nil {
		dbSync.logger.LogError("Failed to delete department, error: ", err)
	}
	return err
}

func (h *QcDepartment) UpdateName(dbSync *DBSync, name string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Name = name
	err := dbSync.UpdateQcDepartment(h)
	return err
}

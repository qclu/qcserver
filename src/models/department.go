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
	hospital, err := dbSync.GetQcHospital(hname)
	if err != nil {
		dbSync.logger.LogError("Failed to delete hospital[", name, "], error: ", err)
		return err
	}
	err = hospital.Delete(dbSync)
	return err
}

func (h *QcHospital) Delete(dbSync *DBSync) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	err := dbSync.DeleteQcHospital(h)
	if err != nil {
		dbSync.logger.LogError("Failed to delete hospital, error: ", err)
	}
	return err
}

func (h *QcHospital) UpdateAddr(dbSync *DBSync, addr string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Addr = addr
	err := dbSync.UpdateQcHospital(h)
	return err
}

func (h *QcHospital) UpdateGis(dbSync *DBSync, gis string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Gis = gis
	return dbSync.UpdateQcHospital(h)
}

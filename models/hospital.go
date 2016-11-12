package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcHospital struct {
	Id      int64      `orm: "pk;auto"`
	Name    string     `orm:"size(256);unique"`
	Addr    string     `orm:"size(2048)"`
	Gis     string     `orm:"size(128)"`
	Created string     `orm:"size(20)"`
	Updated string     `orm:"size(20)"`
	mutex   sync.Mutex `orm:"-"`
}

func CreateQcHospital(dbSync *DBSync, name, addr, gis string) (*QcHospital, error) {
	h := &QcHospital{
		Name:    name,
		Addr:    addr,
		Gis:     gis,
		Created: time.Now().Format(TIME_FMT),
		Updated: time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcHospital(h)
	if err != nil {
		dbSync.logger.LogError("Failed to add new hospital info to database, error: ", err)
		return nil, err
	}
	return h, nil
}

func DeleteQcHospital(dbSync *DBSync, name string) error {
	hospital, err := dbSync.GetQcHospital(name)
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

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcHospital struct {
	Id      int64      `orm: "pk;auto"`
	Name    string     `orm:"size(256);unique"`
	Level   string     `orm:"size(64)"`
	Prov    string     `orm:"size(2048)"`
	City    string     `orm:"size(2048)"`
	Addr    string     `orm:"size(2048)"`
	Gis     string     `orm:"size(128)"`
	Created string     `orm:"size(20)"`
	Updated string     `orm:"size(20)"`
	mutex   sync.Mutex `orm:"-"`
}

func CreateQcHospital(dbSync *DBSync, name, prov, city, addr, gis, level string) (*QcHospital, error) {
	h := &QcHospital{
		Name:    name,
		Prov:    prov,
		Level:   level,
		City:    city,
		Addr:    addr,
		Gis:     gis,
		Created: time.Now().Format(TIME_FMT),
		Updated: time.Now().Format(TIME_FMT),
	}
	dbSync.logger.LogInfo("Hospital to create: ", h)
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

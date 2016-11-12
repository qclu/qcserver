package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcHwVersion struct {
	Id         int64       `orm: "pk;auto"`
	DevModel   *QcDevModel `orm:"rel(fk);on_delete(do_nothing)"`
	Version    string      `orm:"size(256);unique"`
	Annotation string      `orm:"size(4096)"`
	Created    string      `orm:"size(20)"`
	Updated    string      `orm:"size(20)"`
	mutex      sync.Mutex  `orm:"-"`
}

func CreateQcHwVersion(dbSync *DBSync, devmodel *QcDevModel, version string, anno string) (*QcHwVersion, error) {
	obj := &QcHwVersion{
		DevModel:   devmodel,
		Version:    version,
		Annotation: anno,
		Created:    time.Now().Format(TIME_FMT),
		Updated:    time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcHwVersion(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new software version info to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcHwVersion(dbSync *DBSync, version string) error {
	obj, err := dbSync.GetQcHwVersion(version)
	if err != nil {
		dbSync.logger.LogError("Failed to delete software version[", version, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcHwVersion) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcHwVersion(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete software version, error: ", err)
	}
	return err
}

func (obj *QcHwVersion) UpdateVersion(dbSync *DBSync, version string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	obj.Version = version
	err := dbSync.UpdateQcHwVersion(obj)
	return err
}

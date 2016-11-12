package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcSwVersion struct {
	Id          int64        `orm: "pk;auto"`
	HwVersion   *QcHwVersion `orm:"rel(fk);on_delete(do_nothing)"`
	Version     string       `orm:"size(256);unique"`
	SwType      string       `orm:"size(4096)"` //DEBUG | RELEASE
	Description string       `orm:"size(4096)"`
	Created     string       `orm:"size(20)"`
	Updated     string       `orm:"size(20)"`
	mutex       sync.Mutex   `orm:"-"`
}

func CreateQcSwVersion(dbSync *DBSync, hwversion *QcHwVersion, version string, swtype string, description string) (*QcSwVersion, error) {
	obj := &QcSwVersion{
		HwVersion:   hwversion,
		Version:     version,
		SwType:      swtype,
		Description: description,
		Created:     time.Now().Format(TIME_FMT),
		Updated:     time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcSwVersion(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new software version info to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcSwVersion(dbSync *DBSync, version string) error {
	obj, err := dbSync.GetQcSwVersion(version)
	if err != nil {
		dbSync.logger.LogError("Failed to delete software version[", version, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcSwVersion) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcSwVersion(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete software version, error: ", err)
	}
	return err
}

func (obj *QcSwVersion) UpdateVersion(dbSync *DBSync, version string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	obj.Version = version
	err := dbSync.UpdateQcSwVersion(obj)
	return err
}

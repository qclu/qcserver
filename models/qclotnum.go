package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcQcLotnum struct {
	Id          int64        `orm: "pk;auto"`
	Qcp         *QcQcProduct `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Sd          float64      `orm:"digits(12);decimals(4)"`
	TargetValue float64      `orm:"digits(12);decimals(4)"`
	ExpiredTime string       `orm:"size(20);"`
	Type        string       `orm:"size(64)"`
	Created     string       `orm:"size(20)"`
	Updated     string       `orm:"size(20)"`
	mutex       sync.Mutex   `orm:"-"`
}

func CreateQcQcLotnum(dbSync *DBSync, qcp *QcQcProduct, sd, targetval float64, exptime, _type string) (*QcQcLotnum, error) {
	obj := &QcQcLotnum{
		Qcp:         qcp,
		Sd:          sd,
		TargetValue: targetval,
		ExpiredTime: exptime,
		Type:        _type,
		Created:     time.Now().Format(TIME_FMT),
		Updated:     time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcQcLotnum(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new qc lotnum to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcQcLotnum(dbSync *DBSync, id int64) error {
	obj, err := dbSync.GetQcQcLotnum(id)
	if err != nil {
		dbSync.logger.LogError("Failed to delete quality control product lotnum[", id, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcQcLotnum) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcQcLotnum(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete quality control product, error: ", err)
	}
	return err
}

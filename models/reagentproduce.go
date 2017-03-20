package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcReagentProduce struct {
	Id        int64  `orm: "pk;auto"`
	SerialNum string `orm:"size(256);unique"`
	//LotNum      string          `orm:"size(256);unique"`
	RegModel    *QcReagentModel `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	ExpiredTime string          `orm:"size(20)"`
	Annotation  string          `orm:"size(4096)"`
	Created     string          `orm:"size(20)"`
	Updated     string          `orm:"size(20)"`
	mutex       sync.Mutex      `orm:"-"`
}

func CreateQcReagentProduce(dbSync *DBSync, serial, exptime, anno string, regmodel *QcReagentModel) (*QcReagentProduce, error) {
	obj := &QcReagentProduce{
		SerialNum: serial,
		//LotNum:      lognum,
		RegModel:    regmodel,
		ExpiredTime: exptime,
		Annotation:  anno,
		Created:     time.Now().Format(TIME_FMT),
		Updated:     time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcReagentProduce(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new reagent release to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcReagentProduce(dbSync *DBSync, serial string) error {
	obj, err := dbSync.GetQcReagentProduce(serial)
	if err != nil {
		dbSync.logger.LogError("Failed to delete reagent release[", serial, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcReagentProduce) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcReagentProduce(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete reagent release, error: ", err)
	}
	return err
}

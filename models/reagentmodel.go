package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcReagentModel struct {
	Id         int64       `orm: "pk;auto"`
	Name       string      `orm:"size(256);unique"`
	DevModel   *QcDevModel `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Annotation string      `orm:"size(4096)"`
	Created    string      `orm:"size(20)"`
	Updated    string      `orm:"size(20)"`
	mutex      sync.Mutex  `orm:"-"`
}

func CreateQcReagentModel(dbSync *DBSync, name, anno string, devmodel *QcDevModel) (*QcReagentModel, error) {
	obj := &QcReagentModel{
		Name:       name,
		DevModel:   devmodel,
		Annotation: anno,
		Created:    time.Now().Format(TIME_FMT),
		Updated:    time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcReagentModel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new reagent model to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcReagentModel(dbSync *DBSync, name string) error {
	obj, err := dbSync.GetQcReagentModel(name)
	if err != nil {
		dbSync.logger.LogError("Failed to delete reagent model[", name, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcReagentModel) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcReagentModel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete reagent model, error: ", err)
	}
	return err
}

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcDevRel struct {
	Id        int64        `orm: "pk;auto"`
	SwVersion *QcSwVersion `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Sn        string       `orm:"size(256);unique"`
	Date      string       `orm:"size(10)"`
	SmCard    string       `orm:"size(128)"`
	Created   string       `orm:"size(20)"`
	Updated   string       `orm:"size(20)"`
	mutex     sync.Mutex   `orm:"-"`
}

func CreateQcDevModel(dbSync *DBSync, name, model, release string, meth *QcMethodology, anno string) (*QcDevModel, error) {
	obj := &QcDevModel{
		Name:        name,
		Model:       model,
		Release:     release,
		Methodology: meth,
		Annotation:  anno,
		Created:     time.Now().Format(TIME_FMT),
		Updated:     time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcDevModel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new device model to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcDevModel(dbSync *DBSync, name string) error {
	obj, err := dbSync.GetQcDevModel(name)
	if err != nil {
		dbSync.logger.LogError("Failed to delete device model[", name, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcDevModel) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcDevModel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete device model, error: ", err)
	}
	return err
}

func (obj *QcDevModel) UpdateName(dbSync *DBSync, name string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	obj.Name = name
	err := dbSync.UpdateQcDevModel(obj)
	return err
}

func (obj *QcDevModel) UpdateModel(dbSync *DBSync, model string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	obj.Model = model
	err := dbSync.UpdateQcDevModel(obj)
	return err
}

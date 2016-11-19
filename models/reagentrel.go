package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcReagentRel struct {
	Id            int64             `orm: "pk;auto"`
	ReleaseTime   string            `orm:"size(20)"`
	ReleaseSerial string            `orm:"size(256);unique"`
	Department    *QcDepartment     `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	ProduceSerial *QcReagentProduce `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Amounts       float64           `orm:"digits(12);decimals(4)"`
	Annotation    string            `orm:"size(4096)"`
	Consumption   float64           `orm:"digits(12);decimals(4)"`
	Created       string            `orm:"size(20)"`
	Updated       string            `orm:"size(20)"`
	mutex         sync.Mutex        `orm:"-"`
}

func CreateQcReagentRel(dbSync *DBSync, rel_time, rel_serial, anno string, amounts float64, proserial *QcReagentProduce, department *QcDepartment) (*QcReagentRel, error) {
	obj := &QcReagentRel{
		ReleaseTime:   rel_time,
		ReleaseSerial: rel_serial,
		Department:    department,
		ProduceSerial: proserial,
		Amounts:       amounts,
		Annotation:    anno,
		Consumption:   0.0,
		Created:       time.Now().Format(TIME_FMT),
		Updated:       time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcReagentRel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new reagent release to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcReagentRel(dbSync *DBSync, serial string) error {
	obj, err := dbSync.GetQcReagentRel(serial)
	if err != nil {
		dbSync.logger.LogError("Failed to delete reagent release[", serial, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcReagentRel) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcReagentRel(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete reagent release, error: ", err)
	}
	return err
}

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcQcProduct struct {
	Id             int64           `orm: "pk;auto"`
	RegModel       *QcReagentModel `orm:"rel(fk);on_delete(do_nothing)"` // RelForeignKey relation
	Tea            float64         `orm:"digits(12);decimals(4)"`
	Cv             float64         `orm:"digits(12);decimals(4)"`
	Percent        float64         `orm:"digits(12);decimals(4)"`
	FixedDeviation float64         `orm:"digits(12);decimals(4)"`
	Nsd            float64         `orm:"digits(12);decimals(4)"`
	Range          string          `orm:"size(128);decimals(4)"`
	Name           string          `orm:"size(256);unique"`
	Annotation     string          `orm:"size(4096)"`
	Created        string          `orm:"size(20)"`
	Updated        string          `orm:"size(20)"`
	mutex          sync.Mutex      `orm:"-"`
}

func CreateQcQcProduct(dbSync *DBSync, regmodel *QcReagentModel, tea, cv, percent, fixdeviation, nsd float64, _range, name, anno string) (*QcQcProduct, error) {
	obj := &QcQcProduct{
		RegModel:       regmodel,
		Tea:            tea,
		Cv:             cv,
		Percent:        percent,
		FixedDeviation: fixdeviation,
		Nsd:            nsd,
		Range:          _range,
		Name:           name,
		Annotation:     anno,
		Created:        time.Now().Format(TIME_FMT),
		Updated:        time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcQcProduct(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to add new qc product to database, error: ", err)
		return nil, err
	}
	return obj, nil
}

func DeleteQcQcQcProduct(dbSync *DBSync, name string) error {
	obj, err := dbSync.GetQcQcProduct(name)
	if err != nil {
		dbSync.logger.LogError("Failed to delete quality control product[", name, "], error: ", err)
		return err
	}
	err = obj.Delete(dbSync)
	return err
}

func (obj *QcQcProduct) Delete(dbSync *DBSync) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	err := dbSync.DeleteQcQcQcProduct(obj)
	if err != nil {
		dbSync.logger.LogError("Failed to delete quality control product, error: ", err)
	}
	return err
}

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcMethodology struct {
	Id         int64      `orm: "pk;auto"`
	Name       string     `orm:"size(256);unique"`
	Annotation string     `orm:"size(4096)"`
	Created    string     `orm:"size(20)"`
	Updated    string     `orm:"size(20)"`
	mutex      sync.Mutex `orm:"-"`
}

func CreateQcMethodology(dbSync *DBSync, name, anno string) (*QcMethodology, error) {
	h := &QcMethodology{
		Name:       name,
		Annotation: anno,
		Created:    time.Now().Format(TIME_FMT),
		Updated:    time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcMethodology(h)
	if err != nil {
		dbSync.logger.LogError("Failed to add new methodology info to database, error: ", err)
		return nil, err
	}
	return h, nil
}

func DeleteQcMethdology(dbSync *DBSync, name string) error {
	m, err := dbSync.GetQcMethodology(name)
	if err != nil {
		dbSync.logger.LogError("Failed to delete methodology[", name, "], error: ", err)
		return err
	}
	err = m.Delete(dbSync)
	return err
}

func (m *QcMethodology) Delete(dbSync *DBSync) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	err := dbSync.DeleteQcMethdology(m)
	if err != nil {
		dbSync.logger.LogError("Failed to delete methodology, error: ", err)
	}
	return err
}

func (m *QcMethodology) UpdateName(dbSync *DBSync, name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Name = name
	err := dbSync.UpdateQcMethodology(m)
	return err
}

func (m *QcMethodology) UpdateAnnotation(dbSync *DBSync, anno string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Annotation = anno
	return dbSync.UpdateQcMethodology(m)
}

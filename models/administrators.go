package models

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type QcAdministrator struct {
	Id       int64 `orm: "pk;auto"`
	Role     int
	Username string     `orm:"size(64);unique"`
	Password string     `orm:"size(64)"`
	Created  string     `orm:"size(20)"`
	Updated  string     `orm:"size(20)"`
	mutex    sync.Mutex `orm:"-"`
}

func CreateQcAdmin(dbSync *DBSync, username, passwd string, role int) (*QcAdministrator, error) {
	if len(username) <= 0 {
		dbSync.logger.LogError("username cannot be empty...")
		return nil, errors.New("empty username to create administrator")
	}
	if len(passwd) <= 0 {
		dbSync.logger.LogError("password cannot be empty...")
		return nil, errors.New("empty password to create administrator")
	}
	admin := &QcAdministrator{
		Username: username,
		Role:     role,
		Password: passwd,
		Created:  time.Now().Format(TIME_FMT),
		Updated:  time.Now().Format(TIME_FMT),
	}
	err := dbSync.InsertQcAdmin(admin)
	if err != nil {
		dbSync.logger.LogError("Failed to add new administrator info to database, error: ", err)
		return nil, err
	}
	return admin, nil
}

func DeleteQcAdmin(dbSync *DBSync, username string) error {
	admin, err := dbSync.GetQcAdmin(username)
	if err != nil {
		dbSync.logger.LogError("Failed to delete admin[", username, "], error: ", err)
		return err
	}
	err = admin.Delete(dbSync)
	return err
}

func (ad *QcAdministrator) Delete(dbSync *DBSync) error {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()
	err := dbSync.DeleteQcAdmin(ad)
	if err != nil {
		dbSync.logger.LogError("Failed to delete admin, error: ", err)
	}
	return err
}

func (ad *QcAdministrator) UpdatePasswd(dbSync *DBSync, passwd string) error {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()
	ad.Password = passwd
	err := dbSync.UpdateQcAdmin(ad)
	return err
}

func (ad *QcAdministrator) UpdateRole(dbSync *DBSync, role int) error {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()
	ad.Role = role
	return dbSync.UpdateQcAdmin(ad)
}

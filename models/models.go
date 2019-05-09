package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/go-ignite/ignite/config"
)

var (
	db *gorm.DB
)

func MustInitDB() {
	switch config.C.DB.Driver {
	case "mysql", "sqlite3":
	default:
		logrus.WithField("driver", config.C.DB.Driver).Fatal("driver is invalid")
	}
	var err error
	if db, err = gorm.Open(config.C.DB.Driver, config.C.DB.Connect); err != nil {
		logrus.WithError(err).Fatal("connect db error")
	}

	if err := db.DB().Ping(); err != nil {
		logrus.WithError(err).Fatal("ping db error")
	}

	if err := db.AutoMigrate(new(User), new(InviteCode), new(Node), new(Service)).Error; err != nil {
		logrus.WithError(err).Fatal("auto migrate db error")
	}
}

func runTx(tx *gorm.DB, f func() error) error {
	if err := f(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

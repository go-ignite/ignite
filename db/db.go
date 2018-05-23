package db

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var (
	engine *xorm.Engine
	once   sync.Once
)

func GetDB(driver, connect string) *xorm.Engine {
	//Init DB connection
	once.Do(func() {
		switch driver {
		case "mysql", "sqlite3":
		default:
			logrus.WithField("driver", driver).Fatal("driver is invalid")
		}
		var err error
		if engine, err = xorm.NewEngine(driver, connect); err != nil {
			logrus.WithField("err", err).Fatal("new engine error")
		}

		if err := engine.Ping(); err != nil {
			logrus.WithField("err", err).Fatal("engine Ping error")
		}

		if err := engine.Sync2(new(User), new(InviteCode)); err != nil {
			logrus.WithField("err", err).Fatal("models sync error")
		}
	})
	return engine
}

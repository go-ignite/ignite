package db

import (
	"sync"

	"github.com/go-ignite/ignite/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var (
	engine *xorm.Engine
	once   sync.Once
)

func GetDB() *xorm.Engine {
	//Init DB connection
	once.Do(func() {
		switch config.C.DB.Driver {
		case "mysql", "sqlite3":
		default:
			logrus.WithField("driver", config.C.DB.Driver).Fatal("driver is invalid")
		}
		var err error
		if engine, err = xorm.NewEngine(config.C.DB.Driver, config.C.DB.Connect); err != nil {
			logrus.WithField("err", err).Fatal("new engine error")
		}

		if err := engine.Ping(); err != nil {
			logrus.WithField("err", err).Fatal("engine Ping error")
		}

		if err := engine.Sync2(new(User), new(InviteCode), new(Node), new(Service)); err != nil {
			logrus.WithField("err", err).Fatal("models sync error")
		}
	})
	return engine
}

package db

import (
	"database/sql"
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

// Engine represents a xorm engine or session.
type Engine interface {
	Table(tableNameOrBean interface{}) *xorm.Session
	Count(...interface{}) (int64, error)
	Decr(column string, arg ...interface{}) *xorm.Session
	Delete(interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	ID(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Incr(column string, arg ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *xorm.Session
	SQL(interface{}, ...interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
}

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

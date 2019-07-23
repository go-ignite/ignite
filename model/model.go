package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/go-ignite/ignite/config"
)

type Handler struct {
	db *gorm.DB
}

func InitHandler(config config.Model) (*Handler, error) {
	switch config.Driver {
	case "mysql", "sqlite3":
	default:
		return nil, errors.Errorf("model: driver is invalid: %s", config.Driver)
	}

	db, err := gorm.Open(config.Driver, config.Connect)
	if err != nil {
		return nil, errors.Wrap(err, "model: connect to db failed")
	}

	if config.Debug {
		db.LogMode(true)
	}

	if err := db.AutoMigrate(new(User), new(InviteCode), new(Node), new(Service)).Error; err != nil {
		return nil, errors.Wrap(err, "model: db migration error")
	}

	return &Handler{db: db}, nil
}

func (h *Handler) runTX(f func(h *Handler) error) error {
	tx := h.db.Begin()
	if err := f(&Handler{db: tx}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

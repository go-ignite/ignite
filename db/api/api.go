package api

import (
	"github.com/go-ignite/ignite/db"
	"github.com/go-xorm/xorm"
)

type API struct {
	engine  *xorm.Engine
	session *xorm.Session
}

func NewAPI(session ...*xorm.Session) *API {
	api := &API{
		engine: db.GetDB(),
	}
	if len(session) > 0 {
		api.session = session[0]
	}
	return api
}

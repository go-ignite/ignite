package api

import (
	"github.com/go-ignite/ignite/db"
)

type API struct {
	db.Engine
}

func NewAPI(engine ...db.Engine) *API {
	var e db.Engine
	if len(engine) > 0 {
		e = engine[0]
	} else {
		e = db.GetDB()
	}
	return &API{
		Engine: e,
	}
}

package api

import (
	"github.com/go-ignite/ignite/db"
)

func (api *API) GetUserByUsername(username string) (*db.User, error) {
	user := new(db.User)
	_, err := api.Where("username = ?", username).Get(user)
	return user, err
}

func (api *API) GetUserByID(id int64) (*db.User, error) {
	user := new(db.User)
	_, err := api.ID(id).Get(user)
	return user, err
}

func (api *API) UpdateUser(user *db.User, fields ...string) (int64, error) {
	e := api.ID(user.Id)
	if len(fields) > 0 {
		e = e.Cols(fields...)
	}
	return e.Update(user)
}

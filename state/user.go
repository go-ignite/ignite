package state

import (
	"sync"

	"github.com/go-ignite/ignite/model"
)

type User struct {
	sync.RWMutex
	user *model.User
}

func newUser(u *model.User) *User {
	return &User{
		user: u,
	}
}

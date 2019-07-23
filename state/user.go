package state

import (
	"sync"

	"github.com/go-ignite/ignite/model"
)

type User struct {
	user   *model.User
	locker sync.RWMutex
}

func newUser(u *model.User) *User {
	return &User{
		user: u,
	}
}

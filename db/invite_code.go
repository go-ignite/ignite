package db

import (
	"time"
)

type InviteCode struct {
	Id             int64     `xorm:"pk autoincr notnull"`
	InviteCode     string    `xorm:"not null"`
	PackageLimit   int       `xorm:"not null"`
	Available      bool      `xorm:"not null default 1"`
	AvailableLimit int       `xorm:"not null default 1"` //unit: month
	Created        time.Time `xorm:"created"`
	Updated        time.Time `xorm:"updated"`
}

func GetAvailableInviteCode(inviteCode string) (*InviteCode, error) {
	iv := new(InviteCode)
	_, err := engine.Where("invite_code = ? AND available = 1", inviteCode).Get(iv)
	return iv, err
}

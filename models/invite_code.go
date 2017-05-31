package models

import "time"

type InviteCode struct {
	Id           int64     `xorm:"pk autoincr notnull"`
	InviteCode   string    `xorm:"not null"`
	PackageLimit int       `xorm:"not null"`
	Available    bool      `xorm:"not null default 1"`
	UserId       int64     `xorm:"not null default 0"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

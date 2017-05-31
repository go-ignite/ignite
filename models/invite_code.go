package models

import "time"

type InviteCode struct {
	Id           int64     `xorm:"pk autoincr notnull"`
	InviteCode   string    `xorm:"not null"`
	PackageLimit int       `xorm:"not null"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

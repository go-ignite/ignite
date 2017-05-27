package models

import "time"

type User struct {
	Id           int64 `xorm:"pk autoincr notnull"`
	Username     string
	HashedPwd    []byte `xorm:"blob"`
	InviteCode   string
	PackageLimit int `xorm:"not null"`
	PackageUsed  float32
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

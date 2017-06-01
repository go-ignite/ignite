package models

import "time"

type User struct {
	Id           int64 `xorm:"pk autoincr notnull"`
	Username     string
	HashedPwd    []byte `xorm:"blob"`
	InviteCode   string
	PackageLimit int       `xorm:"not null"` //Package bandwidth limit, unit: GB
	PackageUsed  float32   //Package bandwidth used, unit: byte
	ServiceName  string    //Docker service name for SS
	ServicePort  int       `xorm:"not null default 0"` //Docker service port for SS
	ServicePwd   string    //Password for SS
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

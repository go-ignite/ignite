package models

import "time"

type User struct {
	Id           int64 `xorm:"pk autoincr notnull"`
	Username     string
	HashedPwd    []byte `xorm:"blob"`
	InviteCode   string
	PackageLimit int       `xorm:"not null"` //Package bandwidth limit, unit: GB
	PackageUsed  float32   //Package bandwidth used, unit: byte
	Status       int       `xorm:"default 0"`          // 0=>not created 1=>running 2=>stopped
	ServicePort  int       `xorm:"not null default 0"` //Docker service port for SS
	ServicePwd   string    //Password for SS
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
	Expired      time.Time `xorm:"expired"`
}

type UserInfo struct {
	Id                 int64
	Username           string
	Status             int
	PackageLimit       int
	PackageUsed        float32
	PackageUsedPercent string
	ServicePort        int
	ServicePwd         string
}

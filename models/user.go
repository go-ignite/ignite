package models

import "time"

type User struct {
	Id              int64 `xorm:"pk autoincr notnull"`
	Username        string
	HashedPwd       []byte `xorm:"blob"`
	InviteCode      string
	PackageLimit    int        `xorm:"not null"` //Package bandwidth limit, unit: GB
	PackageUsed     float32    //Package bandwidth used, unit: GB
	Status          int        `xorm:"default 0"` // 0=>not created 1=>running 2=>stopped
	ServiceId       string     //SS container id
	ServicePort     int        `xorm:"not null default 0"` //Docker service port for SS
	ServicePwd      string     //Password for SS
	LastStatsResult uint64     //Last time stats result,unit: byte
	LastStatsTime   *time.Time //Last time stats time
	Created         time.Time  `xorm:"created"`
	Updated         time.Time  `xorm:"updated"`
	Expired         time.Time  `xorm:"expired"`
}

type UserInfo struct {
	Id                 int64
	Host               string
	Username           string
	Status             int
	PackageLimit       int
	PackageUsed        string
	PackageLeft        string
	PackageLeftPercent string
	ServicePort        int
	ServicePwd         string
	Expired            string
}

package db

import (
	"time"
)

type Service struct {
	Id              int64 `xorm:"pk autoincr notnull"`
	ServiceID       string
	UserId          int64      `xorm:"notnull"`
	NodeId          int64      `xorm:"notnull"`
	Type            int        // SS/SSR
	Port            int        `xorm:"default 0"` //Docker service port for SS
	Password        string     `xorm:"default ''"`
	Method          string     `xorm:"default ''"`
	Status          int        `xorm:"default 0"`
	LastStatsResult uint64     //Last time stats result,unit: byte
	LastStatsTime   *time.Time //Last time stats time
	Created         time.Time  `xorm:"created"`
	Updated         time.Time  `xorm:"updated"`
}

package db

import (
	"time"
)

type Node struct {
	Id         int64     `xorm:"pk autoincr notnull"`
	Name       string    `xorm:"default '' unique"`
	Comment    string    `xorm:"default ''"`
	Address    string    `xorm:"default ''"`
	Connection string    `xorm:"default ''"`
	PortFrom   int       `xorm:"default 0"`
	PortTo     int       `xorm:"default 0"`
	Services   int       `xorm:"default 0"` // Number of running containers
	Bandwidth  float32   // total bandwidth used (for all containers), unit: GB
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}

package db

import (
	"time"
)

type Node struct {
	Id        int64     `xorm:"pk autoincr notnull"`
	Name      string    `xorm:"not null"`
	Address   string    `xorm:"not null"`
	Available bool      `xorm:"not null default 1"`
	Services  int       `xorm:"default 0"` // Number of running containers
	Bandwidth float32   // total bandwidth used (for all containers), unit: GB
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"` // latest online time
}

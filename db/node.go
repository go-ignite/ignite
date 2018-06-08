package db

import (
	"time"
)

type Node struct {
	Id        int64     `xorm:"pk autoincr notnull"`
	Name      string    `xorm:"default '' unique"`
	Comment   string    `xorm:"default ''"`
	Address   string    `xorm:"default ''"`
	ConnectIP string    `xorm:"default ''"`
	PortFrom  int       `xorm:"default 0"`
	PortTo    int       `xorm:"default 0"`
	Services  int       `xorm:"default 0"` // Number of running containers
	Bandwidth float32   // total bandwidth used (for all containers), unit: GB
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}

func GetAllNodes() ([]*Node, error) {
	nodes := []*Node{}
	return nodes, engine.Find(&nodes)
}

func UpsertNode(node *Node) (int64, error) {
	if node.Id == 0 {
		return engine.Insert(node)
	}
	return engine.Id(node.Id).Cols("name", "comment").Update(node)
}

func DeleteNode(id int64) (int64, error) {
	return engine.Id(id).Delete(new(Node))
}

package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

var ErrDuplicateNodeName = errors.New("node name already exists")

type Node struct {
	ID                string `gorm:"primary_key"`
	Name              string `gorm:"type:varchar(20);unique_index"`
	Comment           string `gorm:"type:varchar(100)"`
	RequestAddress    string `gorm:"type:varchar(20)"`
	ConnectionAddress string `gorm:"type:varchar(20)"`
	PortFrom          uint
	PortTo            uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
}

func NewNode(name, comment, requestAddress, connectionAddress string, portFrom, portTo uint) *Node {
	return &Node{
		ID:                shortuuid.New(),
		Name:              name,
		Comment:           comment,
		RequestAddress:    requestAddress,
		ConnectionAddress: connectionAddress,
		PortFrom:          portFrom,
		PortTo:            portTo,
	}
}

func GetAllNodes() ([]*Node, error) {
	var nodes []*Node
	return nodes, db.Find(&nodes).Error
}

func (n *Node) Save() error {
	return db.Omit("request_address").Save(n).Error
}

func (n *Node) Create() error {
	tx := db.Begin()
	return runTx(tx, func() error {
		if n.checkIfNameExist(tx) {
			return ErrDuplicateNodeName
		}
		return tx.Create(n).Error
	})
}

func CheckIfNodeNameExist(name string) bool {
	n := &Node{Name: name}
	return n.checkIfNameExist(db)
}

func (n *Node) checkIfNameExist(session *gorm.DB) bool {
	count := 0
	session.Model(new(Node)).Where("name = ?", n.Name).Count(&count)
	return count > 0
}

func GetNodeByID(id string) (*Node, error) {
	node := new(Node)
	r := db.First(node, id)
	if r.RecordNotFound() {
		return nil, nil
	}
	return node, r.Error
}

func DeleteNodeByID(id string) error {
	return db.Delete(new(Node), "id = ?", id).Error
}

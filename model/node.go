package model

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
	PortFrom          int
	PortTo            int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
}

func NewNode(name, comment, requestAddress, connectionAddress string, portFrom, portTo int) *Node {
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

func (h *Handler) GetAllNodes() ([]*Node, error) {
	var nodes []*Node
	if err := h.db.Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

func (h *Handler) SaveNode(n *Node) error {
	return h.db.Omit("request_address").Save(n).Error
}

func (h *Handler) CreateNode(n *Node) error {
	return h.runTX(func(tx *gorm.DB) error {
		if n.checkIfNameExist(tx) {
			return ErrDuplicateNodeName
		}
		return tx.Create(n).Error
	})
}

func (h *Handler) CheckIfNodeNameExist(name string) bool {
	n := &Node{Name: name}
	return n.checkIfNameExist(h.db)
}

func (n *Node) checkIfNameExist(db *gorm.DB) bool {
	count := 0
	db.Model(new(Node)).Where("name = ?", n.Name).Count(&count)
	return count > 0
}

func (h *Handler) GetNode(id string) (*Node, error) {
	node := new(Node)
	r := h.db.First(node, id)
	if r.RecordNotFound() {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}

	return node, nil
}

func (h *Handler) DeleteNode(id string) error {
	return h.db.Delete(&Node{ID: id}).Error
}

package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"

	"github.com/go-ignite/ignite/api"
)

type Node struct {
	ID                string `gorm:"primary_key"`
	Name              string `gorm:"type:varchar(20)"`
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

func (n Node) Output() *api.Node {
	return &api.Node{
		ID:                n.ID,
		Name:              n.Name,
		Comment:           n.Comment,
		RequestAddress:    n.RequestAddress,
		ConnectionAddress: n.ConnectionAddress,
		PortFrom:          n.PortFrom,
		PortTo:            n.PortTo,
		CreatedAt:         n.CreatedAt,
	}
}

func (h *Handler) GetAllNodes() ([]*Node, error) {
	var nodes []*Node
	if err := h.db.Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

func (h *Handler) UpdateNode(n *Node) error {
	return h.db.Model(n).Updates(map[string]interface{}{
		"name":               n.Name,
		"comment":            n.Comment,
		"connection_address": n.ConnectionAddress,
		"port_from":          n.PortFrom,
		"port_to":            n.PortTo,
	}).Error
}

func (h *Handler) CreateNode(n *Node) error {
	return h.db.Create(n).Error
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
	r := h.db.First(node, "id = ?", id)
	if r.RecordNotFound() {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}

	return node, nil
}

func (h *Handler) DeleteNode(id string, f func() error) error {
	return h.runTX(func(tx *gorm.DB) error {
		node, err := newHandler(tx).GetNode(id)
		if err != nil {
			return err
		}

		if node == nil {
			return nil
		}

		if err := tx.Delete(Node{ID: id}).Error; err != nil {
			return err
		}

		return f()
	})
}

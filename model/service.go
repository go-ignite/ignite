package model

import (
	"time"

	"github.com/go-ignite/ignite-agent/protos"
)

type ServiceConfig struct {
	Port     uint
	Type     protos.ServiceType_Enum
	Password string
}

type Service struct {
	ID              int64 `gorm:"primary_key"`
	UserID          string
	NodeID          string
	ContainerID     string
	Port            int
	Config          ServiceConfig
	Status          int
	LastStatsResult uint64     // last time stats result,unit: byte
	LastStatsTime   *time.Time // last time stats time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index"`
}

func (h *Handler) GetServices() ([]*Service, error) {
	var services []*Service
	return services, h.db.Find(&services).Error
}

func (h *Handler) GetService(id, userID int64) (*Service, error) {
	r := h.db.Where("id = ?", id)
	if userID > 0 {
		r = r.Where("user_id = ?", userID)
	}

	service := new(Service)
	r = r.First(service)
	if r.RecordNotFound() {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}

	return service, nil
}

func (h *Handler) GetNodePortRangeServiceCount(nodeID string, portFrom, portTo int) (int, error) {
	var count int
	return count, h.db.Model(new(Service)).Where("node_id = ? AND (port < ? OR port > ?)", nodeID, portFrom, portTo).Count(&count).Error
}

func (h *Handler) CreateService(s *Service) error {
	return h.db.Create(s).Error
}

func (h *Handler) GetServicesByUserID(userID int64) ([]*Service, error) {
	var services []*Service
	return services, h.db.Find(&services, "user_id = ?", userID).Error
}

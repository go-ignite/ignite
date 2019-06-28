package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/api"
)

var (
	ErrServiceExists = errors.New("model: user already has a service on this node")
)

type ServiceConfig struct {
	EncryptionMethod protos.ServiceEncryptionMethod_Enum `json:"encryption_method"`
	Password         string                              `json:"password"`
}

func (s ServiceConfig) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ServiceConfig) Scan(value interface{}) error {
	var data []byte
	switch st := value.(type) {
	case []byte:
		data = st
	case string:
		data = []byte(st)
	default:
		return fmt.Errorf("scan src type not matched, get %T", value)
	}

	return json.Unmarshal(data, s)
}

type Service struct {
	ID              int64 `gorm:"primary_key"`
	UserID          string
	NodeID          string
	ContainerID     string
	Type            protos.ServiceType_Enum
	Port            int
	Config          *ServiceConfig `gorm:"type:varchar(1024)"`
	LastStatsResult uint64         // last time stats result,unit: byte
	LastStatsTime   *time.Time     // last time stats time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index"`
}

func NewService(userID, nodeID string, ty protos.ServiceType_Enum, sc *ServiceConfig) *Service {
	return &Service{
		UserID: userID,
		NodeID: nodeID,
		Type:   ty,
		Config: sc,
	}
}

func (s Service) Output() *api.Service {
	return &api.Service{
		ID:               s.ID,
		UserID:           s.UserID,
		NodeID:           s.NodeID,
		Type:             s.Type,
		Port:             s.Port,
		EncryptionMethod: s.Config.EncryptionMethod,
		Password:         s.Config.Password,
		CreatedAt:        s.CreatedAt,
	}
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

func (h *Handler) CreateService(s *Service, f func() error) error {
	return h.runTX(func(tx *gorm.DB) error {
		th := newHandler(tx)
		u, err := th.mustGetUserByID(s.UserID)
		if err != nil {
			return err
		}

		// check if the user create service repeatedly
		count, err := th.checkService(s.UserID, s.NodeID)
		if err != nil {
			return err
		}

		if count > 0 {
			return ErrServiceExists
		}

		s.Config.Password = u.ServicePassword

		// create container so that we can get the port
		if err := f(); err != nil {
			return err
		}

		// TODO there is a problem, create container success but commit transaction error, we need to clean it up
		return tx.Create(s).Error
	})
}

func (h *Handler) checkService(userID, nodeID string) (int, error) {
	var count int
	if err := h.db.Model(Service{}).Where("user_id = ? AND node_id = ?", userID, nodeID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil

}

func (h *Handler) GetServicesByUserID(userID int64) ([]*Service, error) {
	var services []*Service
	return services, h.db.Find(&services, "user_id = ?", userID).Error
}

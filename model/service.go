package model

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/api"
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

func (s Service) Output(host string) *api.Service {
	return &api.Service{
		ID:               s.ID,
		UserID:           s.UserID,
		NodeID:           s.NodeID,
		Type:             s.Type,
		Port:             s.Port,
		EncryptionMethod: s.Config.EncryptionMethod,
		Password:         s.Config.Password,
		CreatedAt:        s.CreatedAt,
		URL:              s.URL(host),
	}
}

func (s Service) URL(host string) string {
	base64Encode := func(s string) string {
		return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(s)), "=")
	}

	var protocol, base64Link string
	switch s.Type {
	case protos.ServiceType_SS_LIBEV:
		protocol = "ss"
		//method:password@server:port
		base64Link = base64Encode(fmt.Sprintf("%s:%s@%s:%d", s.Config.EncryptionMethod.ValidMethod(), s.Config.Password, host, s.Port))
	case protos.ServiceType_SSR:
		protocol = "ssr"
		//server:port:protocol:method:obfs:password_base64/?suffix_base64
		suffix := fmt.Sprintf("protoparam=%s", base64Encode("32"))
		base64Link = base64Encode(fmt.Sprintf("%s:%d:%s:%s:%s:%s/?%s", host, s.Port, "auth_aes128_md5", s.Config.EncryptionMethod.ValidMethod(), "tls1.2_ticket_auth_compatible", base64Encode(s.Config.Password), suffix))
	default:
		return ""
	}
	return fmt.Sprintf("%s://%s", protocol, base64Link)

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

func (h *Handler) CreateService(s *Service) error {
	return h.db.Create(s).Error
}

func (h *Handler) GetServicesByUserID(userID int64) ([]*Service, error) {
	var services []*Service
	return services, h.db.Find(&services, "user_id = ?", userID).Error
}

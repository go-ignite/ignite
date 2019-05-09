package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-ignite/ignite-agent/protos"
)

type Services []*Service

func (services Services) GetNodeIDServicesMap() map[string]Services {
	nodeIDServicesMap := make(map[string]Services)
	for _, service := range services {
		nodeIDServicesMap[service.NodeID] = append(nodeIDServicesMap[service.NodeID], service)
	}

	return nodeIDServicesMap
}

func (services Services) GetNodeIDs() []string {
	var nodeIDs []string
	for _, service := range services {
		nodeIDs = append(nodeIDs, service.NodeID)
	}
	return nodeIDs
}

func (services Services) GetPortMap() map[uint]bool {
	portMap := make(map[uint]bool)
	for _, port := range services.GetPorts() {
		portMap[port] = true
	}

	return portMap
}

func (services Services) GetPorts() []uint {
	var ports []uint
	for _, service := range services {
		ports = append(ports, service.Port)
	}
	return ports
}

func (services Services) IDs() []uint {
	var ids []uint
	for _, service := range services {
		ids = append(ids, service.ID)
	}
	return ids
}

type ServiceConfig struct {
	Port     uint
	Type     protos.ServiceType_Enum
	Password string
	// Method   protos.MethodType_Enum
}

type Service struct {
	gorm.Model
	UserID          uint
	NodeID          string
	ContainerID     string
	Port            uint
	Config          ServiceConfig
	Status          int
	LastStatsResult uint64     // last time stats result,unit: byte
	LastStatsTime   *time.Time // last time stats time
}

func GetAllServices() (Services, error) {
	var services Services
	return services, db.Find(&services).Error
}

func GetServiceByIDAndUserID(id, userID int64) (*Service, error) {
	r := db.Where("id = ?", id)
	if userID > 0 {
		r = r.Where("user_id = ?", userID)
	}

	service := new(Service)
	r = r.First(service)
	if r.RecordNotFound() {
		return nil, nil
	}
	return service, r.Error
}

func GetServiceCountByNodeIDAndPortRange(nodeID string, portFrom, portTo uint) (int, error) {
	var count int
	return count, db.Model(new(Service)).Where("node_id = ? AND (port < ? OR port > ?)", nodeID, portFrom, portTo).Count(&count).Error
}

//func GetServicesByUserIDAndNodeID(userID, nodeID int64) ([]*Service, error) {
//	var session *xorm.Session
//	if userID != 0 {
//		session = engine.Where("user_id = ?", userID)
//	}
//	if nodeID != 0 {
//		session = session.Where("node_id = ?", nodeID)
//	}
//
//	var (
//		services []*Service
//		err      error
//	)
//	if session == nil {
//		err = session.Find(&services)
//	} else {
//		err = session.Find(&services)
//	}
//	return services, err
//}

func deleteServicesByIDs(session *gorm.DB, ids []uint) error {
	return session.Delete(new(Service), "id in (?)", ids).Error
}

func (s *Service) Create() error {
	return db.Create(s).Error
}

//func CheckServiceExists(userID, nodeID int64) (bool, error) {
//	service := &Service{}
//	return engine.Where("user_id = ? AND node_id = ?", userID, nodeID).Get(service)
//}

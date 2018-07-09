package api

import (
	"github.com/go-ignite/ignite/db"
	"github.com/go-xorm/xorm"
)

func (api *API) GetAllServices() ([]*db.Service, error) {
	var services []*db.Service
	return services, api.Find(&services)
}

func (api *API) GetServiceByIDAndUserID(id, userID int64) (*db.Service, error) {
	q := api.Where("id = ?", id)
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	service := new(db.Service)
	_, err := q.Get(service)
	return service, err
}

func (api *API) GetServiceCountByNodeIDAndPortRange(nodeID int64,portFrom,portTo int)(int64,error){
	return api.Where("node_id = ? AND (port < ? OR port > ?)",nodeID,portFrom,portTo).Count(new(db.Service))
}

func (api *API) GetServicesByUserIDAndNodeID(userID, nodeID int64) ([]*db.Service, error) {
	var session *xorm.Session
	if userID != 0 {
		session = api.Where("user_id = ?", userID)
	}
	if nodeID != 0 {
		session = api.Where("node_id = ?", nodeID)
	}

	var (
		services []*db.Service
		err      error
	)
	if session == nil {
		err = api.Find(&services)
	} else {
		err = session.Find(&services)
	}
	return services, err
}

func (api *API) RemoveServiceByID(id int64) (int64, error) {
	return api.ID(id).Delete(new(db.Service))
}

func (api *API) CreateService(service *db.Service) (int64, error) {
	return api.Insert(service)
}

func (api *API) CheckServiceExists(userID, nodeID int64) (bool, error) {
	service := &db.Service{}
	return api.Where("user_id = ? AND node_id = ?", userID, nodeID).Get(service)
}

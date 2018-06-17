package api

import (
	"github.com/go-ignite/ignite/db"
)

func (api *API) GetAllServices() ([]*db.Service, error) {
	var services []*db.Service
	return services, api.engine.Find(&services)
}

func (api *API) GetServiceByIDAndUserIDAndNodeID(id, userID, nodeID int64) (*db.Service, error) {
	service := new(db.Service)
	_, err := api.engine.Where("id = ? AND user_id = ? AND node_id = ?", id, userID, nodeID).Get(service)
	return service, err
}

func (api *API) GetServicesByUserID(userID int64) ([]*db.Service, error) {
	var services []*db.Service
	return services, api.engine.Where("user_id = ?", userID).Find(&services)
}

func (api *API) RemoveServiceByID(id int64) (int64, error) {
	return api.engine.Id(id).Delete(new(db.Service))
}

func (api *API) CreateService(service *db.Service) (int64, error) {
	return api.engine.Insert(service)
}

func (api *API) CheckServiceExists(userID, nodeID int64) (bool, error) {
	service := &db.Service{}
	return api.engine.Where("user_id = ? AND node_id = ?", userID, nodeID).Get(service)
}

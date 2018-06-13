package api

import (
	"github.com/go-ignite/ignite/db"
)

func (api *API) GetAllServices() ([]*db.Service, error) {
	var services []*db.Service
	return services, api.engine.Find(&services)
}

func (api *API) CreateService(service *db.Service) (int64, error) {
	return api.engine.Insert(service)
}

func (api *API) CheckServiceExists(userID, nodeID int64) (bool, error) {
	service := &db.Service{}
	return api.engine.Where("user_id = ? AND node_id = ?", userID, nodeID).Get(service)
}

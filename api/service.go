package api

import (
	"time"

	"github.com/go-ignite/ignite-agent/protos"
)

type CreateServiceRequest struct {
	Type             protos.ServiceType_Enum             `json:"type"`
	EncryptionMethod protos.ServiceEncryptionMethod_Enum `json:"encryption_method"`
	NodeID           string                              `json:"node_id"`
}

type RemoveServiceRequest struct {
	NodeID string `json:"node_id"`
}

type ServiceOptions struct {
	Type              protos.ServiceType_Enum               `json:"type"`
	EncryptionMethods []protos.ServiceEncryptionMethod_Enum `json:"encryption_methods"`
}

type AdminServicesRequest struct {
	UserID string `form:"user_id"`
	NodeID string `form:"node_id"`
}

type Service struct {
	ID               int64                               `json:"id"`
	UserID           string                              `json:"user_id"`
	NodeID           string                              `json:"node_id"`
	Type             protos.ServiceType_Enum             `json:"type"`
	Port             int                                 `json:"port"`
	EncryptionMethod protos.ServiceEncryptionMethod_Enum `json:"encryption_method"`
	Password         string                              `json:"password"`
	CreatedAt        time.Time                           `json:"created_at"`
	URL              string                              `json:"url"`
}

type NodeSyncResponse struct {
	ID        string `json:"id"`
	Available bool   `json:"available"`
}

type ServiceSyncResponse struct {
	ID               int64                     `json:"id"`
	Status           protos.ServiceStatus_Enum `json:"status"`
	MonthTrafficUsed uint64                    `json:"month_traffic_used"`
	LastStatsTime    time.Time                 `json:"last_stats_time"`
}

type NodeServiceSyncResponse struct {
	Node    NodeSyncResponse     `json:"node"`
	Service *ServiceSyncResponse `json:"service"`
}

type UserSyncResponse struct {
	UserID           string                     `json:"user_id"`
	MonthTrafficUsed uint64                     `json:"month_traffic_used"`
	NodeService      []*NodeServiceSyncResponse `json:"node_service"`
}

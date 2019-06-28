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

type Service struct {
	ID               int64                               `json:"id"`
	UserID           string                              `json:"user_id"`
	NodeID           string                              `json:"node_id"`
	Type             protos.ServiceType_Enum             `json:"type"`
	Port             int                                 `json:"port"`
	EncryptionMethod protos.ServiceEncryptionMethod_Enum `json:"encryption_method"`
	Password         string                              `json:"password"`
	CreatedAt        time.Time                           `json:"created_at"`
}

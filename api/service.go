package api

import "github.com/go-ignite/ignite-agent/protos"

type ServiceConfig struct {
	Type    protos.ServiceType_Enum `json:"type"`
	Methods []string                `json:"methods"`
}

type CreateServiceRequest struct {
	Type             protos.ServiceType_Enum             `json:"type"`
	EncryptionMethod protos.ServiceEncryptionMethod_Enum `json:"encryption_method"`
	NodeID           string                              `json:"node_id"`
}

type ServiceInfoResp struct {
	Id       int64  `json:"id"`
	NodeId   int64  `json:"node_id"`
	Type     int    `json:"type"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Status   int    `json:"status"`
	Created  int64  `json:"created"`
}

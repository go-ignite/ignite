package api

import "github.com/go-ignite/ignite-agent/protos"

type ServiceConfig struct {
	Type    protos.ServiceType_Enum `json:"type"`
	Methods []string                `json:"methods"`
}

type CreateServiceReq struct {
	Type     protos.ServiceType_Enum `json:"type"`
	Method   string                  `json:"method"`
	Password string                  `json:"password"`
	NodeID   int64                   `json:"-"`
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

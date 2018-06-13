package models

import pb "github.com/go-ignite/ignite-agent/protos"

type ServiceConfig struct {
	Type      string              `json:"type"`
	TypeProto pb.ServiceType_Enum `json:"-"`
	Methods   []string            `json:"methods"`
}

type CreateServiceReq struct {
	Type     string `json:"type"`
	Method   string `json:"method"`
	Password string `json:"password"`
	NodeID   int64  `json:"-"`
}

type CreateServiceResp struct {
	Type     string `json:"type"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

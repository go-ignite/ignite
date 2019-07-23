package state

import (
	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/model"
)

type Service struct {
	service *model.Service
	user    *User
	status  protos.ServiceStatus_Enum
}

func newService(s *model.Service, u *User) *Service {
	return &Service{
		service: s,
		user:    u,
	}
}

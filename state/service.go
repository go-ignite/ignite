package state

import (
	"sync"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/model"
)

type Service struct {
	lock    sync.RWMutex
	service *model.Service
	status  protos.ServiceStatus_Enum
}

func newService(s *model.Service) *Service {
	return &Service{
		service: s,
	}
}

func (s *Service) updateSyncResponse(resp *protos.ServiceInfo) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// TODO
}

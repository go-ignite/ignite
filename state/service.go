package state

import "github.com/go-ignite/ignite/model"

type Service struct {
	service   *model.Service
	available bool
}

func newService(s *model.Service) *Service {
	return &Service{
		service: s,
	}
}

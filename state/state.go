package state

import (
	"context"
	"sync"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/model"
)

var Set = wire.NewSet(wire.Struct(new(Options), "*"), Init)

type Handler struct {
	nodes sync.Map
	opts  *Options
}

type Options struct {
	Config       config.State
	ModelHandler *model.Handler
}

func Init(opts *Options) (*Handler, error) {
	h := &Handler{
		opts: opts,
	}

	nodes, err := h.opts.ModelHandler.GetAllNodes()
	if err != nil {
		return nil, err
	}

	services, err := h.opts.ModelHandler.GetServices()
	if err != nil {
		return nil, err
	}

	nodeServices := map[string][]*model.Service{}
	for _, s := range services {
		nodeServices[s.NodeID] = append(nodeServices[s.NodeID], s)
	}
	for _, node := range nodes {
		n, err := newNode(h.opts.Config, node, nodeServices[node.ID])
		if err != nil {
			return nil, err
		}
		h.nodes.Store(node.ID, n)
	}

	return h, nil
}

func (h *Handler) Start() {
	h.nodes.Range(func(_, n interface{}) bool {
		go n.(*Node).monitor()
		return true
	})
}

func (h *Handler) AddNode(ctx context.Context, node *model.Node) error {
	var err error
	h.nodes.Range(func(_, v interface{}) bool {
		n := v.(*Node)
		switch {
		case n.node.Name == node.Name:
			err = api.ErrNodeNameExists
		case n.node.RequestAddress == node.RequestAddress:
			err = api.ErrNodeRequestAddressExists
		default:
			return true
		}

		return false
	})
	if err != nil {
		return err
	}

	n, err := newNode(h.opts.Config, node, nil)
	if err != nil {
		return err
	}

	if _, err := grpc_health_v1.NewHealthClient(n.conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: protos.ServiceName,
	}); err != nil {
		return err
	}

	if err := h.opts.ModelHandler.CreateNode(node); err != nil {
		return err
	}

	go n.monitor()
	h.nodes.Store(n.node.ID, n)

	return nil
}

func (h *Handler) UpdateNode(node *model.Node) error {
	n1, ok := h.nodes.Load(node.ID)
	if !ok {
		return api.ErrNodeNotExist
	}

	n := n1.(*Node)
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.node.PortFrom != node.PortFrom || n.node.PortTo != node.PortTo {
		count := 0
		for _, s := range n.services {
			if s.service.Port < node.PortFrom || s.service.Port > node.PortTo {
				count++
			}
		}
		if count > 0 {
			return api.ErrNodeHasServicesExceedPortRange
		}
	}

	if err := h.opts.ModelHandler.UpdateNode(node); err != nil {
		return err
	}

	n.node.Name = node.Name
	n.node.Comment = node.Comment
	n.node.ConnectionAddress = node.ConnectionAddress
	n.node.PortFrom = node.PortFrom
	n.node.PortTo = node.PortTo
	return nil
}

func (h *Handler) RemoveNode(nodeID string) error {
	n1, ok := h.nodes.Load(nodeID)
	if !ok {
		return nil
	}

	if err := h.opts.ModelHandler.DeleteNode(nodeID); err != nil {
		return err
	}

	// TODO clean up node containers
	n1.(*Node).stopMonitor()
	h.nodes.Delete(nodeID)
	return nil
}

func (h *Handler) AddService(service *model.Service) error {
	n1, ok := h.nodes.Load(service.NodeID)
	if !ok {
		return api.ErrNodeNotExist
	}

	n := n1.(*Node)
	if !n.isAvailable() {
		return api.ErrNodeUnavailable
	}

	if n.checkServiceExist(service.UserID) {
		return api.ErrServiceExists
	}

	req := &protos.CreateServiceRequest{
		PortFrom:         int32(n.node.PortFrom),
		PortTo:           int32(n.node.PortTo),
		Type:             service.Type,
		EncryptionMethod: service.Config.EncryptionMethod,
		Password:         service.Config.Password,
		UserId:           service.UserID,
		NodeId:           service.NodeID,
	}

	resp, err := n.client.CreateService(context.Background(), req)
	if err != nil {
		// TODO distinguish error
		return err
	}
	service.Port = int(resp.Port)
	service.ContainerID = resp.ContainerId

	if err := h.opts.ModelHandler.CreateService(service); err != nil {
		// failed to create service, clean it up
		removeReq := &protos.RemoveServiceRequest{
			ContainerId: resp.ContainerId,
		}
		if _, err := n.client.RemoveService(context.Background(), removeReq); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"containerID": resp.ContainerId,
				"userID":      service.UserID,
			}).Error("create container successfully, but save to db error, then cleaning the container failed")
		}
		return err
	}

	n.addService(service)
	return nil
}

func (h *Handler) GetNodeServices(userID, nodeID string) []*api.NodeServices {
	nss := make([]*api.NodeServices, 0)

	getUserServices := func(n *Node) *api.NodeServices {
		n.lock.RLock()
		defer n.lock.RUnlock()

		ns := &api.NodeServices{}
		ns.Node = n.node.Output()

		if userID != "" {
			if s, ok := n.services[userID]; ok {
				ns.Services = append(ns.Services, s.service.Output(n.node.ConnectionAddress))
			}
		} else {
			for _, s := range n.services {
				ns.Services = append(ns.Services, s.service.Output(n.node.ConnectionAddress))
			}
		}

		return ns
	}

	if nodeID != "" {
		n1, ok := h.nodes.Load(nodeID)
		if !ok {
			return nss
		}

		nss = append(nss, getUserServices(n1.(*Node)))
		return nss
	}

	h.nodes.Range(func(_, n1 interface{}) bool {
		nss = append(nss, getUserServices(n1.(*Node)))
		return true
	})

	return nss
}

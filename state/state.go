package state

import (
	"context"
	"errors"
	"sync"

	"github.com/google/wire"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/model"
)

var Set = wire.NewSet(wire.Struct(new(Options), "*"), Init)

var (
	ErrNodeNotExist                   = errors.New("state: node does not exist")
	ErrNodeUnavailable                = errors.New("state: node is unavailable")
	ErrNodeHasServicesExceedPortRange = errors.New("state: node has services that exceed port range")
)

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
	n, err := newNode(h.opts.Config, node, nil)
	if err != nil {
		return err
	}

	if _, err := n.client.Ping(ctx, new(protos.PingRequest)); err != nil {
		return err
	}

	go n.monitor()
	h.nodes.Store(n.node.ID, n)

	return nil
}

func (h *Handler) UpdateNode(node *model.Node, f func() error) error {
	n1, ok := h.nodes.Load(node.ID)
	if !ok {
		return ErrNodeNotExist
	}
	if err := f(); err != nil {
		return err
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
			return ErrNodeHasServicesExceedPortRange
		}
	}

	n.node.Name = node.Name
	n.node.Comment = node.Comment
	n.node.ConnectionAddress = node.ConnectionAddress
	n.node.PortFrom = node.PortFrom
	n.node.PortTo = node.PortTo
	return nil
}

func (h *Handler) RemoveNode(nodeID string) {
	n, ok := h.nodes.Load(nodeID)
	if !ok {
		return
	}

	n.(*Node).stopMonitor()
	h.nodes.Delete(nodeID)
}

func (h *Handler) AddService(ctx context.Context, service *model.Service) error {
	n1, ok := h.nodes.Load(service.NodeID)
	if !ok {
		return ErrNodeNotExist
	}

	n := n1.(*Node)
	if !n.isAvailable() {
		return ErrNodeUnavailable
	}

	req := &protos.CreateServiceRequest{
		PortFrom:         int32(n.node.PortFrom),
		PortTo:           int32(n.node.PortTo),
		Type:             service.Type,
		EncryptionMethod: service.Config.EncryptionMethod,
		Password:         service.Config.Password,
		Name:             service.UserID,
		NodeId:           service.NodeID,
	}

	resp, err := n.client.CreateService(ctx, req)
	if err != nil {
		// TODO distinguish error
		return err
	}

	service.Port = int(resp.Port)
	n.addService(service)

	return nil
}

func (h *Handler) GetUserServices(userID string) []*api.NodeService {
	nss := make([]*api.NodeService, 0)
	h.nodes.Range(func(_, n interface{}) bool {
		ns := &api.NodeService{}
		node := n.(*Node)
		ns.Node = node.node.Output()
		if s, ok := node.services[userID]; ok {
			ns.Service = s.service.Output()

		}

		nss = append(nss, ns)
		return true
	})

	return nss
}

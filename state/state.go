package state

import (
	"context"
	"sync"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/model"
)

var Set = wire.NewSet(
	wire.Struct(new(Options), "*"),
	Init,
)

type Handler struct {
	nodes       map[string]*Node
	nodesLocker sync.RWMutex
	users       map[string]*User
	usersLocker sync.RWMutex
	opts        *Options
}

type Options struct {
	Config       config.State
	ModelHandler *model.Handler
}

func Init(opts *Options) (*Handler, error) {
	h := &Handler{
		opts:  opts,
		users: map[string]*User{},
		nodes: map[string]*Node{},
	}

	users, err := h.opts.ModelHandler.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		h.users[user.ID] = newUser(user)
	}

	nodes, err := h.opts.ModelHandler.GetAllNodes()
	if err != nil {
		return nil, err
	}

	services, err := h.opts.ModelHandler.GetServices()
	if err != nil {
		return nil, err
	}

	nodeServices := map[string][]*Service{}
	for _, s := range services {
		nodeServices[s.NodeID] = append(nodeServices[s.NodeID], newService(s, h.users[s.UserID]))
	}

	h.nodesLocker.Lock()
	defer h.nodesLocker.Unlock()

	for _, node := range nodes {
		n, err := newNode(h.opts.Config, node, nodeServices[node.ID], opts.ModelHandler)
		if err != nil {
			return nil, err
		}

		h.nodes[node.ID] = n
	}

	return h, nil
}

func (h *Handler) Start() {
	h.nodesLocker.RLock()
	defer h.nodesLocker.RUnlock()

	for _, node := range h.nodes {
		go node.monitor()
	}

	go h.runDailyTask()
}

func (h *Handler) AddNode(ctx context.Context, n *model.Node) error {
	h.nodesLocker.Lock()
	defer h.nodesLocker.Unlock()

	for _, node := range h.nodes {
		var err error
		switch {
		case node.node.Name == n.Name:
			err = api.ErrNodeNameExists
		case node.node.RequestAddress == n.RequestAddress:
			err = api.ErrNodeRequestAddressExists
		}

		if err != nil {
			return err
		}
	}

	node, err := newNode(h.opts.Config, n, nil, h.opts.ModelHandler)
	if err != nil {
		return err
	}

	if _, err := grpc_health_v1.NewHealthClient(node.conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: protos.ServiceName,
	}); err != nil {
		return err
	}

	if err := h.opts.ModelHandler.CreateNode(n); err != nil {
		return err
	}

	go node.monitor()
	h.nodes[node.node.ID] = node

	return nil
}

func (h *Handler) UpdateNode(n *model.Node) error {
	node, unlock := h.getLockedNode(n.ID)
	if node == nil {
		return api.ErrNodeNotExist
	}
	defer unlock()

	if node.node.PortFrom != n.PortFrom || node.node.PortTo != n.PortTo {
		count := 0
		for _, s := range node.services {
			if s.service.Port < n.PortFrom || s.service.Port > n.PortTo {
				count++
			}
		}
		if count > 0 {
			return api.ErrNodeHasServicesExceedPortRange
		}
	}

	if err := h.opts.ModelHandler.UpdateNode(n); err != nil {
		return err
	}

	node.node.Name = n.Name
	node.node.Comment = n.Comment
	node.node.ConnectionAddress = n.ConnectionAddress
	node.node.PortFrom = n.PortFrom
	node.node.PortTo = n.PortTo
	return nil
}

func (h *Handler) RemoveNode(nodeID string) error {
	h.nodesLocker.Lock()
	defer h.nodesLocker.Unlock()

	node, ok := h.nodes[nodeID]
	if !ok {
		return nil
	}

	if err := h.opts.ModelHandler.DeleteNode(nodeID); err != nil {
		return err
	}

	// TODO clean up node containers
	node.stopMonitor()
	delete(h.nodes, nodeID)
	return nil
}

func (h *Handler) AddService(service *model.Service) error {
	node, unlock := h.getLockedNode(service.NodeID)
	if node == nil {
		return api.ErrNodeNotExist
	}
	defer unlock()

	if !node.available {
		return api.ErrNodeUnavailable
	}

	if _, ok := node.services[service.UserID]; ok {
		return api.ErrServiceExists
	}

	req := &protos.CreateServiceRequest{
		PortFrom:         int32(node.node.PortFrom),
		PortTo:           int32(node.node.PortTo),
		Type:             service.Type,
		EncryptionMethod: service.Config.EncryptionMethod,
		Password:         service.Config.Password,
		UserId:           service.UserID,
		NodeId:           service.NodeID,
	}

	resp, err := node.client.CreateService(context.Background(), req)
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
		if _, err := node.client.RemoveService(context.Background(), removeReq); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"containerID": resp.ContainerId,
				"userID":      service.UserID,
			}).Error("create container successfully, but save to db error, then cleaning the container failed")
		}
		return err
	}

	h.usersLocker.RLock()
	defer h.usersLocker.RUnlock()

	node.services[service.UserID] = newService(service, h.users[service.UserID])
	return nil
}

func (h *Handler) GetNodeServices(userID, nodeID string) []*api.NodeServices {
	nss := make([]*api.NodeServices, 0)

	getUserServices := func(n *Node) *api.NodeServices {
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
		node, unlock := h.getRLockedNode(nodeID)
		if node == nil {
			return nss
		}
		defer unlock()

		nss = append(nss, getUserServices(node))
		return nss
	}

	h.nodesLocker.RLock()
	defer h.nodesLocker.RUnlock()

	for _, node := range h.nodes {
		nss = append(nss, getUserServices(node))
	}

	return nss
}

func (h *Handler) GetSyncResponse(userID string) []*api.UserSyncResponse {
	h.usersLocker.RLock()
	defer h.usersLocker.RUnlock()

	r := make([]*api.UserSyncResponse, 0)
	for _, user := range h.users {
		if userID != "" && user.user.ID != userID {
			continue
		}

		func() {
			user.locker.RLock()
			defer user.locker.RUnlock()

			usr := &api.UserSyncResponse{
				UserID: user.user.ID,
			}

			h.nodesLocker.RLock()
			defer h.nodesLocker.RUnlock()

			for _, node := range h.nodes {
				func() {
					node.locker.RLock()
					defer node.locker.RUnlock()

					nodeService := &api.NodeServiceSyncResponse{
						Node: api.NodeSyncResponse{
							ID:        node.node.ID,
							Available: node.available,
						},
					}

					if service := node.services[user.user.ID]; service != nil {
						nodeService.Service = &api.ServiceSyncResponse{
							ID:               service.service.ID,
							Status:           service.status,
							MonthTrafficUsed: service.service.MonthTrafficUsed(),
							LastStatsTime:    service.service.LastStatsTime,
						}
						usr.MonthTrafficUsed += nodeService.Service.MonthTrafficUsed
					}

					usr.NodeService = append(usr.NodeService, nodeService)

				}()
			}

			r = append(r, usr)
		}()
	}

	return r
}

func (h *Handler) AddUser(user *model.User) error {
	if err := h.opts.ModelHandler.CreateUser(user); err != nil {
		return err
	}

	h.usersLocker.Lock()
	defer h.usersLocker.Unlock()

	h.users[user.ID] = newUser(user)
	return nil
}

func (h *Handler) RemoveUser(userID string) error {
	if err := h.opts.ModelHandler.DestroyUser(userID); err != nil {
		return err
	}

	// TODO clean up containers

	h.usersLocker.Lock()
	defer h.usersLocker.Unlock()

	delete(h.users, userID)
	return nil
}

func (h *Handler) CheckUserExists(userID string) bool {
	h.usersLocker.RLock()
	defer h.usersLocker.RUnlock()

	return h.users[userID] != nil
}

func (h *Handler) ChangeUserPassword(userID, newPassword string, oldPassword *string) error {
	h.usersLocker.RLock()
	defer h.usersLocker.RUnlock()

	u := h.users[userID]
	u.locker.Lock()
	defer u.locker.Unlock()

	if oldPassword != nil {
		if err := bcrypt.CompareHashAndPassword(u.user.HashedPwd, []byte(*oldPassword)); err != nil {
			return api.ErrUserPasswordIncorrect
		}
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := h.opts.ModelHandler.ChangeUserPassword(userID, hashedPwd); err != nil {
		return err
	}

	u.user.HashedPwd = hashedPwd
	return nil
}

func (h *Handler) getLockedNode(nodeID string) (*Node, func()) {
	h.nodesLocker.RLock()

	node := h.nodes[nodeID]
	if node == nil {
		h.nodesLocker.RUnlock()
		return nil, nil
	}

	node.locker.Lock()

	f := func() {
		node.locker.Unlock()
		h.nodesLocker.RUnlock()
	}

	return node, f
}

func (h *Handler) getRLockedNode(nodeID string) (*Node, func()) {
	h.nodesLocker.RLock()

	node := h.nodes[nodeID]
	if node == nil {
		h.nodesLocker.RUnlock()
		return nil, nil
	}

	node.locker.RLock()

	f := func() {
		node.locker.RUnlock()
		h.nodesLocker.RUnlock()
	}

	return node, f
}

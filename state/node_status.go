package state

import (
	"context"
	"sync"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/logger"
	"github.com/go-ignite/ignite/models"
)

type NodeStatus struct {
	sync.RWMutex
	node        *models.Node
	usedPortMap map[uint]bool

	watching  bool
	available bool

	logger *logger.Logger
	client *agent.Client
}

func NewNodeStatus(node *models.Node, client *agent.Client, usedPortMap map[uint]bool) *NodeStatus {
	if usedPortMap == nil {
		usedPortMap = map[uint]bool{}
	}
	return &NodeStatus{
		node:        node,
		usedPortMap: usedPortMap,
		logger:      logger.GetAgentLogger(),
		client:      client,
	}
}

func (ns *NodeStatus) Heartbeat() error {
	if ns.client.AgentServiceClient == nil {
		if err := ns.client.Dial(); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	streamClient, err := ns.client.NodeHeartbeat(ctx, new(pb.GeneralRequest))
	if err != nil {
		return err
	}

	for {
		if _, err := streamClient.Recv(); err != nil {
			ns.client.AgentServiceClient = nil
			ns.available = false
			return err
		} else {
			ns.available = true
		}
	}
	return nil
}

func (ns *NodeStatus) Available() bool {
	return ns.available
}

func (ns *NodeStatus) Node() *models.Node {
	return ns.node
}

func (ns *NodeStatus) Client() *agent.Client {
	return ns.client
}

func (ns *NodeStatus) GetUsedPorts() []int32 {
	ns.RLock()
	defer ns.RUnlock()

	var ports []int32
	for port := range ns.usedPortMap {
		ports = append(ports, int32(port))
	}
	return ports
}

func (ns *NodeStatus) AddPortToUsedMap(port uint) {
	ns.Lock()
	defer ns.Lock()

	ns.usedPortMap[port] = true
}

func (ns *NodeStatus) RemovePortFromUsedMap(port uint) {
	ns.Lock()
	defer ns.Unlock()

	delete(ns.usedPortMap, port)
}

package state

import (
	"sync"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type NodeStatus struct {
	Node        *db.Node
	UsedPortMap map[int]bool
	available   bool

	watching bool
	*agent.Client
	sync.RWMutex
	*logrus.Logger
}

func NewNodeStatus(node *db.Node, client *agent.Client, available bool, usedPortMap map[int]bool) *NodeStatus {
	if usedPortMap == nil {
		usedPortMap = map[int]bool{}
	}
	return &NodeStatus{
		Node:        node,
		UsedPortMap: usedPortMap,
		Client:      client,
		available:   available,
	}
}

func (ns *NodeStatus) Heartbeat() error {
	if ns.Client.AgentServiceClient == nil {
		if err := ns.Client.Dial(); err != nil {
			return err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streamClient, err := ns.Client.NodeHeartbeat(ctx, &pb.GeneralRequest{})
	if err != nil {
		return err
	}
	for {
		_, err := streamClient.Recv()
		if err != nil {
			ns.Client.AgentServiceClient = nil
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

func (ns *NodeStatus) GetUsedPorts() []int32 {
	var ports []int32
	for port := range ns.UsedPortMap {
		ports = append(ports, int32(port))
	}
	return ports
}

func (ns *NodeStatus) AddPortToUsedMap(port int) {
	ns.Lock()
	defer ns.Lock()

	ns.UsedPortMap[port] = true
}

func (ns *NodeStatus) RemovePortFromUsedMap(port int) {
	ns.Lock()
	defer ns.Lock()

	delete(ns.UsedPortMap, port)
}

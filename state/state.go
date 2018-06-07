package state

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"

	"github.com/sirupsen/logrus"
)

var (
	l  *Loader
	no sync.Once
)

type Loader struct {
	nodeMutex sync.RWMutex
	nodeMap   map[int64]*NodeStatus
	*logrus.Logger
}

type NodeStatus struct {
	*agent.Client
	watching  bool
	available bool
}

func NewNodeStatus(client *agent.Client, available bool) *NodeStatus {
	return &NodeStatus{
		Client:    client,
		available: available,
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

func GetLoader() *Loader {
	no.Do(func() {
		l = &Loader{
			nodeMap: map[int64]*NodeStatus{},
		}
	})
	return l
}

func (loader *Loader) Load() error {
	loader.nodeMutex.Lock()
	defer loader.nodeMutex.Unlock()

	nodes, err := db.GetAllNodes()
	if err != nil {
		return fmt.Errorf("db.GetAllNodes error: %v", err)
	}
	for _, node := range nodes {
		client := agent.NewClient(node.Address)
		ns := NewNodeStatus(client, false)
		go loader.WatchNode(ns)
		loader.nodeMap[node.Id] = ns
	}
	return nil
}

func (loader *Loader) WatchNode(ns *NodeStatus) {
	ns.watching = true
	for ns.watching {
		if err := ns.Heartbeat(); err != nil {
			loader.WithError(err).Error()
			time.Sleep(5 * time.Second)
		}
	}
}

func (loader *Loader) GetNodeAvailable(id int64) bool {
	loader.nodeMutex.RLock()
	defer loader.nodeMutex.RUnlock()

	node := loader.nodeMap[id]
	if node == nil {
		return false
	}
	return node.available
}

func (loader *Loader) GetNodeAvailableMap() map[int64]bool {
	loader.nodeMutex.RLock()
	defer loader.nodeMutex.RUnlock()

	nam := map[int64]bool{}
	for id, ns := range loader.nodeMap {
		nam[id] = ns.available
	}
	return nam
}

func (loader *Loader) DelNode(id int64) {
	loader.nodeMutex.Lock()
	defer loader.nodeMutex.Unlock()

	ns := loader.nodeMap[id]
	if ns != nil {
		ns.watching = false
		ns.Client.Close()
		delete(loader.nodeMap, id)
	}
}

func (loader *Loader) AddNode(id int64, ns *NodeStatus) {
	loader.nodeMutex.Lock()
	defer loader.nodeMutex.Unlock()

	go loader.WatchNode(ns)
	loader.nodeMap[id] = ns
}

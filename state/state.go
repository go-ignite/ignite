package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db/api"

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

	nodes, err := api.NewAPI().GetAllNodes()
	if err != nil {
		return fmt.Errorf("db.GetAllNodes error: %v", err)
	}
	services, err := api.NewAPI().GetAllServices()
	if err != nil {
		return fmt.Errorf("db.GetAllServices error: %v", err)
	}

	nodePortMap := map[int64]map[int]bool{}
	for _, service := range services {
		portMap, ok := nodePortMap[service.NodeId]
		if !ok {
			portMap = map[int]bool{}
			nodePortMap[service.NodeId] = portMap
		}
		portMap[service.Port] = true
	}

	for _, node := range nodes {
		client := agent.NewClient(node.Address)
		ns := NewNodeStatus(node, client, false, nodePortMap[node.Id])
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

func (loader *Loader) GetNode(id int64) *NodeStatus {
	loader.nodeMutex.Lock()
	defer loader.nodeMutex.Unlock()

	return loader.nodeMap[id]
}

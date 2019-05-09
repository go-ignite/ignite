package state

import (
	"sync"
	"time"

	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/logger"
	"github.com/go-ignite/ignite/models"
)

var loader *Loader

type Loader struct {
	sync.RWMutex
	nodeStatusMap map[string]*NodeStatus
	logger        *logger.Logger
}

func NewLoader() *Loader {
	return &Loader{
		nodeStatusMap: make(map[string]*NodeStatus),
		logger:        logger.GetAgentLogger(),
	}
}

func GetLoader() *Loader {
	return loader
}

func MustLoad() {
	loader = NewLoader()
	nodes, err := models.GetAllNodes()
	if err != nil {
		loader.logger.WithError(err).Fatal("models.GetAllNodes error")
	}
	services, err := models.GetAllServices()
	if err != nil {
		loader.logger.WithError(err).Fatal("models.GetAllServices error")
	}

	nodeIDServicesMap := services.GetNodeIDServicesMap()
	for _, node := range nodes {
		ns := NewNodeStatus(node, agent.NewClient(node.RequestAddress), nodeIDServicesMap[node.ID].GetPortMap())
		go loader.WatchNodeStatus(ns)
		loader.nodeStatusMap[node.ID] = ns
	}
}

func (l *Loader) WatchNodeStatus(ns *NodeStatus) {
	ns.watching = true
	for ns.watching {
		if err := ns.Heartbeat(); err != nil {
			l.logger.WithError(err).WithField("nodeID", ns.node.ID).Error("node status Heartbeat error")
			// TODO should be configurable
			time.Sleep(5 * time.Second)
		}
	}
}

func (l *Loader) NodeIsAvailable(nodeID string) bool {
	l.RLock()
	defer l.RUnlock()

	ns, ok := l.nodeStatusMap[nodeID]
	if !ok {
		return false
	}

	return ns.available
}

func (l *Loader) NodeAvailableMap() map[string]bool {
	l.RLock()
	defer l.RUnlock()

	nam := map[string]bool{}
	for nodeID, ns := range l.nodeStatusMap {
		nam[nodeID] = ns.available
	}

	return nam
}

func (l *Loader) DeleteNodeStatus(id string) {
	l.Lock()
	defer l.Unlock()

	ns := l.nodeStatusMap[id]
	if ns != nil {
		ns.watching = false
		ns.client.Close()
		delete(l.nodeStatusMap, id)
	}
}

func (l *Loader) AddNodeStatus(ns *NodeStatus) {
	l.Lock()
	defer l.Unlock()

	ns.logger = logger.GetAgentLogger()
	l.nodeStatusMap[ns.node.ID] = ns
	go l.WatchNodeStatus(ns)
}

func (l *Loader) GetNodeStatus(nodeID string) *NodeStatus {
	l.Lock()
	defer l.Unlock()

	return l.nodeStatusMap[nodeID]
}

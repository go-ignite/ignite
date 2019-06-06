package state

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/model"
)

type Node struct {
	lock      sync.RWMutex
	node      *model.Node
	services  map[string]*Service
	ports     map[int]bool
	available bool
	conn      *grpc.ClientConn
	client    pb.AgentServiceClient
	done      chan struct{}
}

func newNode(node *model.Node, services []*model.Service) (*Node, error) {
	conn, err := grpc.Dial(node.RequestAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	n := &Node{
		node:     node,
		services: map[string]*Service{},
		ports:    map[int]bool{},
		conn:     conn,
		client:   pb.NewAgentServiceClient(conn),
		done:     make(chan struct{}),
	}
	for _, s := range services {
		n.services[s.UserID] = newService(s)
		n.ports[s.Port] = true
	}

	return n, nil
}

func (ns *Node) sync() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ns.done
		cancel()
	}()

	for {
		if err := func() error {
			streamClient, err := ns.client.NodeHeartbeat(ctx, new(pb.GeneralRequest))
			if err != nil {
				return err
			}

			for {
				if _, err := streamClient.Recv(); err != nil {
					ns.available = false
					return err
				} else {
					ns.available = true
				}
			}
		}(); err != nil {
			if err == context.Canceled {
				break
			}

			logrus.WithError(err).WithField("nodeID", ns.node.ID).Error("state: node sync error")
		}

		time.Sleep(3 * time.Second)
	}

	ns.done <- struct{}{}
}

func (ns *Node) stopSync() {
	ns.done <- struct{}{}
	<-ns.done

	_ = ns.conn.Close()
}

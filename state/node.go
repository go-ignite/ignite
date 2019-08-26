package state

import (
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/model"
)

const GB = 1024 * 1024 * 1024

type token string

func (t token) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": string(t),
	}, nil
}

func (token) RequireTransportSecurity() bool {
	return false
}

type Node struct {
	sync.RWMutex
	modelHandler *model.Handler
	config       config.State
	node         *model.Node
	services     map[string]*Service
	available    bool
	conn         *grpc.ClientConn
	client       protos.AgentServiceClient
	done         chan struct{}
}

func newNode(config config.State, node *model.Node, services []*Service, modelHandler *model.Handler) (*Node, error) {
	conn, err := grpc.Dial(
		node.RequestAddress,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(token(config.AgentToken)),
	)
	if err != nil {
		return nil, err
	}

	n := &Node{
		modelHandler: modelHandler,
		node:         node,
		services:     map[string]*Service{},
		conn:         conn,
		config:       config,
		client:       protos.NewAgentServiceClient(conn),
		done:         make(chan struct{}),
	}
	for _, s := range services {
		n.services[s.service.UserID] = s
	}

	return n, nil
}

func (n *Node) heartbeat(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if err := func() error {
			req := &protos.HeartbeatRequest{
				Interval: ptypes.DurationProto(n.config.HeartbeatInterval),
			}
			stream, err := n.client.Heartbeat(ctx, req)
			if err != nil {
				return err
			}

			for {
				_, err := stream.Recv()
				func() {
					n.Lock()
					defer n.Unlock()

					available := err == nil
					if !available && n.available {
						for _, s := range n.services {
							s.status = protos.ServiceStatus_NOT_SET
						}
					}
					n.available = available
				}()
				if err != nil {
					return err
				}
			}
		}(); err != nil {
			if err == context.Canceled {
				return
			}

			logrus.WithError(err).WithField("nodeID", n.node.ID).Error("state: node heartbeat error")
		}

		select {
		case <-time.After(n.config.StreamRetryInterval):
		case <-ctx.Done():
			return
		}
	}
}

func (n *Node) sync(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if err := func() error {
			req := &protos.SyncRequest{
				SyncInterval: ptypes.DurationProto(n.config.SyncInterval),
				NodeId:       n.node.ID,
			}
			stream, err := n.client.Sync(ctx, req)
			if err != nil {
				return err
			}

			for {
				resp, err := stream.Recv()
				if err != nil {
					return err
				}
				n.applySyncResp(resp)
			}
		}(); err != nil {
			if err == context.Canceled {
				return
			}

			logrus.WithError(err).WithField("nodeID", n.node.ID).Error("state: node sync error")
		}

		select {
		case <-time.After(n.config.StreamRetryInterval):
		case <-ctx.Done():
			return
		}
	}
}

func (n *Node) monitor() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-n.done
		cancel()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go n.heartbeat(ctx, wg)
	go n.sync(ctx, wg)

	wg.Wait()
	close(n.done)
}

func (n *Node) stopMonitor() {
	n.done <- struct{}{}
	<-n.done

	_ = n.conn.Close()
}

func (n *Node) applySyncResp(resp *protos.SyncStreamServer) {
	now := time.Now()
	n.Lock()
	defer n.Unlock()

	for _, s := range resp.Services {
		func() {
			if service, ok := n.services[s.ContainerName]; ok {
				service.status = s.Status
				service.service.LastStatsResult = uint64(s.StatsResult)
				service.service.LastStatsTime = now

				if err := n.modelHandler.UpdateServiceLastStats(service.service); err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"nodeID":    n.node.ID,
						"userID":    service.user.user.ID,
						"serviceID": service.service.ID,
					}).Error("state: update service last stats error")
				}

				if service.user != nil {
					service.user.RLock()
					defer service.user.RUnlock()

					if service.service.MonthTrafficUsed() >= uint64(service.user.user.PackageLimit*GB) && service.status == protos.ServiceStatus_RUNNING {
						req := &protos.StopServiceRequest{
							ContainerId: s.ContainerId,
						}
						if _, err := n.client.StopService(context.Background(), req); err != nil {
							logrus.WithError(err).WithFields(logrus.Fields{
								"nodeID":    n.node.ID,
								"userID":    service.user.user.ID,
								"serviceID": service.service.ID,
							}).Error("state: user exceed package limit, but stop service error")
						}
					}
				}
			}
		}()
	}
}

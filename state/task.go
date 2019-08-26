package state

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-ignite/ignite-agent/protos"
)

func (h *Handler) runDailyTask() {
	go func() {
		now := time.Now()
		nextTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)

		select {
		case <-time.After(nextTime.Sub(now)):
			go h.dailyTask(nextTime)
		}

		ticker := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-ticker.C:
				nextTime = nextTime.AddDate(0, 0, 1)
				go h.dailyTask(nextTime)
			}
		}
	}()
}

func (h *Handler) dailyTask(t time.Time) {
	h.nodesLocker.RLock()
	defer h.nodesLocker.RUnlock()

	h.usersLocker.RLock()
	defer h.usersLocker.RUnlock()

	expiredUserMap := map[string]bool{}
	for _, user := range h.users {
		if user.user.ExpiredAt.Before(t) {
			expiredUserMap[user.user.ID] = true
		}
	}

	for nodeID, node := range h.nodes {
		func() {
			node.Lock()
			defer node.Unlock()

			for _, s := range node.services {
				if expiredUserMap[s.service.UserID] && s.status == protos.ServiceStatus_RUNNING {
					req := &protos.StopServiceRequest{
						ContainerId: s.service.ContainerID,
					}
					if _, err := node.client.StopService(context.Background(), req); err != nil {
						logrus.WithError(err).WithFields(logrus.Fields{
							"nodeID":    nodeID,
							"userID":    s.service.UserID,
							"serviceID": s.service.ID,
						}).Error("state: daily stats task, user expired, but stop service error")
					}
				}

				if s.service.MonthBaseStatsTime.Year() != t.Year() || s.service.MonthBaseStatsTime.Month() != t.Month() {
					if !expiredUserMap[s.service.UserID] && s.status != protos.ServiceStatus_RUNNING {
						req := &protos.StartServiceRequest{
							ContainerId: s.service.ContainerID,
						}
						if _, err := node.client.StartService(context.Background(), req); err != nil {
							logrus.WithError(err).WithFields(logrus.Fields{
								"nodeID":    nodeID,
								"userID":    s.service.UserID,
								"serviceID": s.service.ID,
							}).Error("state: daily stats task, start service error")
						}
					}
				}

				if err := h.opts.ModelHandler.UpdateServiceStats(s.service, t); err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"nodeID":    nodeID,
						"userID":    s.service.UserID,
						"serviceID": s.service.ID,
					}).Error("state: update service stats info failed")
				}
			}
		}()
	}
}

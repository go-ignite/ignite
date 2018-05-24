package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/ss"
)

const (
	GB = 1024 * 1024 * 1024
)

//instantStats: Instant task, check & update used bandwidth, stop containers which exceeded the package limit.
func (ctx *CronJob) InstantStats() {

	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	// 1. Load all service from user
	users := []db.User{}
	err := db.GetDB().Where("service_id != ''").Find(&users)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		return
	}

	// 2. Compute ss bandwidth
	for _, user := range users {
		if ss.IsContainerRunning(user.ServiceId) {
			user.Status = 1
			raw, err := ss.GetContainerStatsOutNet(user.ServiceId)
			if err != nil {
				log.Printf("Get container(%s) net out error: %s\n", user.ServiceId, err.Error())
				continue
			}

			// Get container start time
			startTime, err := ss.GetContainerStartTime(user.ServiceId)
			if err != nil {
				log.Printf("Get container(%s) start time error: %s\n", user.ServiceId, err.Error())
				continue
			}

			// Update user package used
			var bandwidth float32
			if user.LastStatsTime == nil || user.LastStatsTime.Before(*startTime) {
				bandwidth = float32(float64(raw) / GB)
			} else {
				bandwidth = float32(float64(raw-user.LastStatsResult) / GB)
			}
			user.PackageUsed += bandwidth

			if int(user.PackageUsed) >= user.PackageLimit {
				// Stop container && update user status
				err := ss.StopContainer(user.ServiceId)
				if err != nil {
					log.Printf("Stop container(%s) error: %s\n", user.ServiceId, err.Error())
				} else {
					log.Printf("STOP: user(%d-%s)-container(%s)\n", user.Id, user.Username, user.ServiceId[:12])
					user.Status = 2
					user.PackageUsed = float32(user.PackageLimit)
				}
			}
			now := time.Now()
			user.LastStatsResult = raw
			user.LastStatsTime = &now
			if b := fmt.Sprintf("%.2f", bandwidth); b != "0.00" {
				log.Printf("STATS: user(%d-%s)-container(%s)-bandwidth(%s)\n", user.Id, user.Username, user.ServiceId[:12], b)
			}
		} else {
			user.Status = 2
		}

		// 3. Update user stats info
		_, err = db.GetDB().Id(user.Id).Cols("package_used", "last_stats_result", "last_stats_time", "status").Update(user)
		if err != nil {
			log.Printf("Update user(%d) error: %s\n", user.Id, err.Error())
			continue
		}
	}
}

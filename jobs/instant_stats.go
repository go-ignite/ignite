package jobs

import (
	"log"
	"os"
	"time"

	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-xorm/xorm"
)

const (
	GB = 1024 * 1024 * 1024
)

var (
	db *xorm.Engine
)

func SetDB(engine *xorm.Engine) {
	db = engine
}

//instantStats: Instant task, check & update used bandwidth, stop containers which exceeded the package limit.
func InstantStats() {
	// 1. Load all service from user
	users := []models.User{}
	err := db.Where("service_id != '' AND status = 1").Find(&users)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		os.Exit(1)
	}

	// 2. Compute ss bandwidth
	for _, user := range users {
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
				log.Println("Stop container(%s) error: %s\n", user.ServiceId, err.Error())
			} else {
				log.Printf("STOP: user(%d-%s)-container(%s)\n", user.Id, user.Username, user.ServiceId[:12])
				user.Status = 2
				user.PackageUsed = float32(user.PackageLimit)
			}
		}

		// 3. Update user stats info
		now := time.Now()
		user.LastStatsTime = &now
		user.LastStatsResult = raw
		_, err = db.Id(user.Id).Cols("package_used", "last_stats_result", "last_stats_time", "status").Update(user)
		if err != nil {
			log.Printf("Update user(%d) error: %s\n", user.Id, err.Error())
			continue
		}
		log.Printf("STATS: user(%d-%s)-container(%s)-bandwidth(%.2f)\n", user.Id, user.Username, user.ServiceId[:12], bandwidth)
	}
}

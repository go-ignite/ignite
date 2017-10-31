package jobs

import (
	"log"
	"os"
	"time"

	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
)

//dailyStats: Daily task, check & stop expired containers.
func DailyStats() {
	//1. Load all services from users
	users := []models.User{}
	err := db.Where("service_id != '' AND status = 1").Find(&users)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		os.Exit(1)
	}

	//2. Stop expired containers
	for _, user := range users {
		if user.Expired.Before(time.Now()) {
			err = ss.StopContainer(user.ServiceId)

			if err == nil {
				user.Status = 2
				user.PackageUsed = float32(user.PackageLimit)
				_, err = db.Id(user.Id).Cols("package_used", "status").Update(user)
				if err != nil {
					log.Printf("Update user(%d) error: %s\n", user.Id, err.Error())
					continue
				}
				log.Printf("Stop container:%s for user:%s \r\n", user.ServiceId, user.Username)
			}
		}
	}
}

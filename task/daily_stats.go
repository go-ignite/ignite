package task

import (
	"log"
	"time"

	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/ss"
)

//dailyStats: Daily task, check & stop expired containers.
func (t *Task) DailyStats() {
	t.Lock()
	defer t.Unlock()

	//1. Load all services from users
	users := []db.User{}
	err := db.GetDB().Where("service_id != '' AND status = 1").Find(&users)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		return
	}

	//2. Stop expired containers
	for _, user := range users {
		if user.Expired.Before(time.Now()) {
			err = ss.KillContainer(user.ServiceId)

			if err == nil {
				user.Status = 2
				user.PackageUsed = float32(user.PackageLimit)
				_, err = db.GetDB().Id(user.Id).Cols("package_used", "status").Update(user)
				if err != nil {
					log.Printf("Update user(%d) error: %s\n", user.Id, err.Error())
					continue
				}
				log.Printf("Stop container:%s for user:%s \r\n", user.ServiceId, user.Username)
			}
		}
	}
}

package jobs

import (
	"log"
	"time"

	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/ss"
)

//monthlyStats: Restart stopped containers and restore the bandwidth.
func (ctx *CronJob) MonthlyStats() {

	ctx.mux.Lock()
	defer ctx.mux.Unlock()

	log.Println("In MonthlyStats")
	log.Println("Load all stopped services from users")
	//1. Load all stopped services from users
	users := []db.User{}
	err := db.GetDB().Where("service_id != '' AND status = 2").Find(&users)
	if err == nil {
		//2. Restart stopped but not expired containers
		for _, user := range users {
			if user.Expired.After(time.Now()) {
				err = ss.StartContainer(user.ServiceId)

				if err == nil {
					user.Status = 1
					user.PackageUsed = 0
					_, err = db.GetDB().Id(user.Id).Cols("package_used", "status").Update(user)

					if err != nil {
						log.Printf("Update user(%d) error: %s\n", user.Id, err.Error())
						continue
					}

					log.Printf("Start container:%s for user:%s \r\n", user.ServiceId, user.Username)
				}
			}
		}
	} else {
		log.Println("Get users error: ", err.Error())
	}

	log.Println("Start load all running services from users")
	//3. Load all running services from  users
	runningUsers := []db.User{}
	err = db.GetDB().Where("service_id != '' AND status = 1").Find(&runningUsers)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		return
	}

	//4. Reset used package for all running users
	for _, ru := range runningUsers {
		if ru.Expired.After(time.Now()) {
			ru.PackageUsed = 0
			_, err = db.GetDB().Id(ru.Id).Cols("package_used").Update(ru)

			if err != nil {
				log.Printf("Update user(%d) error: %s\n", ru.Id, err.Error())
				continue
			}
		}
	}
}

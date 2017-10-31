package jobs

import (
	"log"
	"os"
	"time"

	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/ss"
)

//monthlyStats: Restart stopped containers and restore the bandwidth.
func MonthlyStats() {
	//1. Load all stopped services from users
	users := []models.User{}
	err := db.Where("service_id != '' AND status = 2").Find(&users)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		os.Exit(1)
	}

	//2. Restart stopped but not expired containers
	for _, user := range users {
		if user.Expired.After(time.Now()) {
			err = ss.StartContainer(user.ServiceId)

			if err == nil {
				user.Status = 1
				user.PackageUsed = 0
				_, err = db.Id(user.Id).Cols("package_used", "status").Update(user)

				if err != nil {
					log.Printf("Update user(%d) error: %s\n", user.Id, err.Error())
					continue
				}

				log.Printf("Start container:%s for user:%s \r\n", user.ServiceId, user.Username)
			}
		}
	}

	//3. Load all running services from  users
	runningUsers := []models.User{}
	err = db.Where("service_id != '' AND status = 1").Find(&runningUsers)
	if err != nil {
		log.Println("Get users error: ", err.Error())
		os.Exit(1)
	}

	//4. Reset used package for all running users
	for _, ru := range runningUsers {
		if ru.Expired.After(time.Now()) {
			ru.PackageUsed = 0
			_, err = db.Id(ru.Id).Cols("package_used").Update(ru)

			if err != nil {
				log.Printf("Update user(%d) error: %s\n", ru.Id, err.Error())
				continue
			}
		}
	}
}

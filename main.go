package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/controllers"
	"github.com/go-ignite/ignite/jobs"
	"github.com/robfig/cron"
)

var (
	versionFlag = flag.Bool("v", false, "version")
	version     = "unknown"
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		return
	}

	config.Init()

	if len(os.Args) > 1 && os.Args[1] == "recover" {
		jobs.RecoverTask()
		return
	}

	go initJob()
	mainRouter := &controllers.MainRouter{}
	mainRouter.Initialize(gin.Default())
}

func initJob() {
	c := cron.New()
	cj := &jobs.CronJob{}

	c.AddFunc("0 */5 * * * *", cj.InstantStats)
	c.AddFunc("0 0 0 * * *", cj.DailyStats)
	c.AddFunc("0 0 0 1 * *", cj.MonthlyStats)
	c.Start()
	select {}
}

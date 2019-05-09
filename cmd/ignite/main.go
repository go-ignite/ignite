package main

import (
	"flag"
	"fmt"

	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/handler"
	"github.com/go-ignite/ignite/logger"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/service"
	"github.com/go-ignite/ignite/state"
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

	// init config
	config.MustInit()

	// init logger
	logger.MustInit()

	// init db
	models.MustInitDB()

	// init loader
	state.MustLoad()

	// start task
	//task.New().AsyncRun()

	// start service
	service.New(
		handler.NewUserHandler(),
		handler.NewAdminHandler(),
	).Init().Run()
}

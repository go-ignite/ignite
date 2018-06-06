package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/handler"
	"github.com/go-ignite/ignite/logger"
	"github.com/go-ignite/ignite/router"
	"github.com/go-ignite/ignite/state"
	"github.com/go-ignite/ignite/task"
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

	// init db
	db.GetDB()

	// start task
	t := task.New(logger.New("task.log"))
	t.Init()
	t.Start()

	// init loader
	loader := state.GetLoader()
	loader.Logger = logger.New("agent.log")
	loader.Load()

	r := &router.Router{
		Engine:       gin.Default(),
		UserHandler:  handler.NewUserHandler(logger.New("user.log")),
		AdminHandler: handler.NewAdminHandler(logger.New("admin.log")),
	}
	r.Init()
	r.Run(config.C.App.Address)
}

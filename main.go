package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/controllers"
	"github.com/go-ignite/ignite/utils"
)

var confPath = flag.String("c", "./conf/config.toml", "config file")

func main() {
	flag.Parse()
	utils.InitConf(*confPath)
	initRouter()
}

func initRouter() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	mainRouter := &controllers.MainRouter{}
	mainRouter.Initialize(r)
}

package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/conf"
	"github.com/go-ignite/ignite/controllers"
)

func main() {
	flag.Parse()
	conf.InitConf()
	initRouter()
}

func initRouter() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	mainRouter := &controllers.MainRouter{}
	mainRouter.Initialize(r)
}

package controllers

import (
	"log"

	"github.com/go-ignite/ignite/config"
	_ "github.com/go-ignite/ignite/docs"
	"github.com/go-ignite/ignite/middleware"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-ignite/ignite/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type MainRouter struct {
	router *gin.Engine
	db     *xorm.Engine
}

func (self *MainRouter) Initialize(r *gin.Engine) {
	ss.Host = config.C.Host.Address
	ss.PortRange = []int{config.C.Host.From, config.C.Host.To}

	self.router = r
	self.db = utils.InitDB(config.C.DB.Driver, config.C.DB.Connect)

	if gin.Mode() == gin.DebugMode {
		self.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := self.router.Group("/api")

	user := api.Group("/user")
	{
		user.POST("/login", self.LoginHandler)
		user.POST("/signup", self.SignupHandler)

		auth := user.Group("/auth")
		auth.Use(middleware.Auth(config.C.Auth.Secret))
		{
			auth.GET("/info", self.PanelIndexHandler)
			auth.GET("/config", self.PanelIndexHandler)
			auth.POST("/create", self.CreateServiceHandler)
		}
	}

	go func() {
		if err := ss.PullImage(ss.SS_IMAGE); err != nil {
			log.Printf("Pull image [%s] error: %s\n", ss.SS_IMAGE, err.Error())
		}
		if err := ss.PullImage(ss.SSR_IMAGE); err != nil {
			log.Printf("Pull image [%s] error: %s\n", ss.SSR_IMAGE, err.Error())
		}
	}()
	self.router.Run(config.C.APP.Address)
}

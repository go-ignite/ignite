package controllers

import (
	"log"

	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/db"
	_ "github.com/go-ignite/ignite/docs"
	"github.com/go-ignite/ignite/middleware"
	"github.com/go-ignite/ignite/ss"

	"github.com/gin-gonic/contrib/cors"
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
	self.db = db.GetDB()

	if gin.Mode() == gin.DebugMode {
		self.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		self.router.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
	}

	api := self.router.Group("/api")

	user := api.Group("/user")
	{
		user.POST("/login", self.LoginHandler)
		user.POST("/signup", self.SignupHandler)

		auth := user.Group("/auth")
		auth.Use(middleware.Auth(config.C.Secret.User))
		{
			auth.GET("/info", self.UserInfoHandler)
			auth.GET("/config", self.ServiceConfigHandler)
			auth.POST("/create", self.CreateServiceHandler)
		}
	}
	admin := api.Group("/admin")
	{
		admin.POST("/login", self.PanelLoginHandler)
		auth := admin.Group("/auth")
		auth.Use(middleware.Auth(config.C.Secret.Admin))
		{
			//user account related operations
			auth.GET("/status_list", self.PanelStatusListHandler)
			auth.PUT("/:id/reset", self.ResetAccountHandler)
			auth.PUT("/:id/destroy", self.DestroyAccountHandler)
			auth.PUT("/:id/stop", self.StopServiceHandler)
			auth.PUT("/:id/start", self.StartServiceHandler)
			auth.PUT("/:id/renew", self.RenewServiceHandler)

			//invite code related operations
			auth.GET("/code_list", self.InviteCodeListHandler)
			auth.PUT("/:id/remove", self.RemoveInviteCodeHandler)
			auth.POST("/code_generate", self.GenerateInviteCodeHandler)
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

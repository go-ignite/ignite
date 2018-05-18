package controllers

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-ignite/ignite/utils"
	"github.com/go-xorm/xorm"
)

type MainRouter struct {
	router *gin.Engine
	db     *xorm.Engine
}

func (self *MainRouter) Initialize(r *gin.Engine) {
	//Init session store
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("ignite", store))
	ss.Host = config.C.Host.Address
	ss.PortRange = []int{config.C.Host.From, config.C.Host.To}

	self.router = r
	self.router.GET("/", self.IndexHandler)
	self.router.POST("/login", self.LoginHandler)
	self.router.POST("/signup", self.SignupHandler)
	self.db = utils.InitDB(config.C.DB.Driver, config.C.DB.Connect)

	pg := self.router.Group("/panel")
	pg.Use(ValidateSession())
	{
		pg.GET("/index", self.PanelIndexHandler)
		pg.GET("/logout", self.LogoutHandler)
		pg.POST("/create", self.CreateServiceHandler)
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

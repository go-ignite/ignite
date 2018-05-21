package controllers

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/ss"
	"github.com/go-ignite/ignite/utils"
	"github.com/go-xorm/xorm"
)

type MainRouter struct {
	router *gin.Engine
	db     *xorm.Engine
}

func (self *MainRouter) Initialize(r *gin.Engine) {
	ss.Host = utils.HOST_Address
	ss.PortRange = []int{utils.HOST_From, utils.HOST_To}

	//Init session store
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("ignite", store))

	self.router = r
	self.db = utils.InitDB(utils.DB_Driver, utils.DB_Connect)
	self.router.GET("/", self.IndexHandler)
	self.router.POST("/login", self.LoginHandler)
	self.router.POST("/signup", self.SignupHandler)

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
	self.router.Run(utils.APP_Address)
}

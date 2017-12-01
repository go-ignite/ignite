package controllers

import (
	"github.com/gin-contrib/sessions"
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
	ss.ImageUrl = utils.SS_Image
	ss.PortRange = []int{utils.HOST_From, utils.HOST_To}

	//Init session store
	store := sessions.NewCookieStore([]byte("secret"))
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

	self.router.Run(utils.APP_Address)
}

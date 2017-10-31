package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/conf"
	"github.com/go-ignite/ignite/ss"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type MainRouter struct {
	router *gin.Engine
	db     *xorm.Engine
}

func (self *MainRouter) Initialize(r *gin.Engine, db *xorm.Engine) {
	ss.Host = conf.HOST_Address
	ss.ImageUrl = conf.SS_Image
	ss.PortRange = []int{conf.HOST_From, conf.HOST_To}

	//Init session store
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("ignite", store))

	self.router = r
	self.db = db
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

	self.router.Run(conf.APP_Address)
}

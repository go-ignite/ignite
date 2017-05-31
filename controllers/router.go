package controllers

import (
	"flag"
	"fmt"
	"ignite/models"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	toml "github.com/pelletier/go-toml"
)

var (
	conf = flag.String("c", "./config.toml", "config file")
)

type MainRouter struct {
	router *gin.Engine
	db     *xorm.Engine
}

func (self *MainRouter) Initialize(r *gin.Engine) {
	//Check config file
	if _, err := os.Stat(*conf); os.IsNotExist(err) {
		fmt.Println("Cannot load config.toml, file doesn't exist...")
		os.Exit(1)
	}

	config, err := toml.LoadFile(*conf)

	if err != nil {
		fmt.Println("Failed to load config file:", *conf)
		os.Exit(1)
	}

	//Init DB connection
	connString := fmt.Sprintf("%s:%s@/ignite?charset=utf8", config.Get("mysql.user").(string), config.Get("mysql.password").(string))
	engine, _ := xorm.NewEngine("mysql", connString)

	err = engine.Ping()

	if err != nil {
		fmt.Println("Cannot connetc to database:", err.Error())
		os.Exit(1)
	}

	err = engine.Sync2(new(models.User), new(models.InviteCode))
	if err != nil {
		fmt.Println("Failed to sync database struct:", err.Error())
		os.Exit(1)
	}

	self.db = engine

	//Init session store
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("ignite", store))

	self.router = r
	self.router.GET("/", self.IndexHandler)
	self.router.POST("/login", self.LoginHandler)
	self.router.POST("/signup", self.SignupHandler)

	pg := self.router.Group("/panel")
	{
		pg.GET("/", self.PanelIndexHandler)
	}

	self.router.Run(":5000")
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/conf"
	"github.com/go-ignite/ignite/controllers"
	"github.com/go-ignite/ignite/jobs"
	"github.com/go-ignite/ignite/models"
	"github.com/go-xorm/xorm"
	"github.com/robfig/cron"
)

func main() {
	flag.Parse()
	conf.InitConf()

	db := initDB()
	go initRouter(db)
	go initJob(db)
	select {}
}

func initRouter(db *xorm.Engine) {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	mainRouter := &controllers.MainRouter{}
	mainRouter.Initialize(r, db)
}

func initJob(db *xorm.Engine) {
	jobs.SetDB(db)
	c := cron.New()
	c.AddFunc("* */5 * * * *", jobs.InstantStats)
	c.AddFunc("0 0 0 * * *", jobs.DailyStats)
	c.AddFunc("0 0 1 * * *", jobs.MonthlyStats)
	c.Start()
	select {}
}

func initDB() *xorm.Engine {
	//Init DB connection
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.MYSQL_User, conf.MYSQL_Password, conf.MYSQL_Host, conf.MYSQL_DBName)
	engine, _ := xorm.NewEngine("mysql", connString)

	err := engine.Ping()

	if err != nil {
		fmt.Println("Cannot connetc to database:", err.Error())
		os.Exit(1)
	}

	err = engine.Sync2(new(models.User), new(models.InviteCode))
	if err != nil {
		fmt.Println("Failed to sync database struct:", err.Error())
		os.Exit(1)
	}
	return engine
}

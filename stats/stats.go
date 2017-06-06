package main

import (
	"flag"
	"fmt"
	"ignite/ss"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-xorm/xorm"
	toml "github.com/pelletier/go-toml"
)

var (
	conf = flag.String("c", "./config.toml", "config file")
	db   *xorm.Engine
)

func init() {
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
	var (
		user     = config.Get("mysql.user").(string)
		password = config.Get("mysql.password").(string)
		host     = config.Get("mysql.host").(string)
		dbname   = config.Get("mysql.dbname").(string)
	)
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, password, host, dbname)
	db, err = xorm.NewEngine("mysql", connString)
	if err != nil {
		fmt.Println("Create mysql engine error:", err.Error())
		os.Exit(1)
	}

	err = db.Ping()

	if err != nil {
		fmt.Println("Cannot connetc to database:", err.Error())
		os.Exit(1)
	}
}

func main() {
	raw, err := ss.StatsOutNet("de80655c7a2d11359e698ea0d0a92d21bc22f076b1b3aa5f0dadfb255594c65e")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(raw)
	//bandwidth = raw / 1024 / 1024

	//if bandwidth > 1024:
	//print('Bandwidth is: {} / {:.2f} GB'.format(raw, float(bandwidth)/1024))
	//else:
	//print('Bandwidth is: {} / {:.2f} MB'.format(raw, float(bandwidth)))
}

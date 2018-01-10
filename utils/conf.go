package utils

import (
	"fmt"
	"os"
	"strconv"

	toml "github.com/pelletier/go-toml"
)

var (
	// for app config
	APP_Address string

	// for db config
	DB_Driver, DB_Connect string

	// for host config
	HOST_Address       string
	HOST_From, HOST_To int
)

func InitConf(confPath string) {
	//Check config file
	if _, err := os.Stat(confPath); !os.IsNotExist(err) {
		if config, err := toml.LoadFile(confPath); err == nil {
			APP_Address = config.Get("app.address").(string)

			HOST_Address = config.Get("host.address").(string)
			HOST_From = int(config.Get("host.from").(int64))
			HOST_To = int(config.Get("host.to").(int64))

			DB_Driver = config.Get("db.driver").(string)
			DB_Connect = config.Get("db.connect").(string)
		}
	}
	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		DB_Driver = driver
	}
	if connect := os.Getenv("DB_CONNECT"); connect != "" {
		DB_Connect = connect
	}
	if address := os.Getenv("HOST_ADDRESS"); address != "" {
		HOST_Address = address
	}
	if from := os.Getenv("HOST_FROM"); from != "" {
		HOST_From, _ = strconv.Atoi(from)
	}
	if to := os.Getenv("HOST_TO"); to != "" {
		HOST_To, _ = strconv.Atoi(to)
	}
	fmt.Println("config: ", map[string]interface{}{
		"address":      APP_Address,
		"db_driver":    DB_Driver,
		"db_connect":   DB_Connect,
		"host_address": HOST_Address,
		"host_from":    HOST_From,
		"host_to":      HOST_To,
	})
}

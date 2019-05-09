package config

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	C Config
)

type (
	Admin struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	Config struct {
		App struct {
			Address string `mapstructure:"address"`
			Secret  string `mapstructure:"secret"`
		} `mapstructure:"app"`
		Log struct {
			Task         string `mapstructure:"task"`
			Agent        string `mapstructure:"agent"`
			UserHandler  string `mapstructure:"user_handler"`
			AdminHandler string `mapstructure:"admin_handler"`
			Level        string `mapstructure:"level"`
		} `mapstructure:"log"`
		DB struct {
			Driver  string `mapstructure:"driver"`
			Connect string `mapstructure:"connect"`
		} `mapstructure:"db"`
		Host struct {
			Address string `mapstructure:"address"`
			From    int    `mapstructure:"from"`
			To      int    `mapstructure:"to"`
		} `mapstructure:"host"`
		Admin Admin `mapstructure:"admin"`
	}
)

func (a *Admin) Match(username, password string) bool {
	return a.Username == username && a.Password == password
}

func (c *Config) MustCheck() {
	if _, err := logrus.ParseLevel(C.Log.Level); err != nil {
		log.Fatalf("parse app.log_level error: %v\n", err)
	}
}

func MustInit() {
	// app
	viper.SetDefault("app.address", ":5000")
	viper.SetDefault("app.secret", "ignite")
	// log
	viper.SetDefault("log.task", "task.log")
	viper.SetDefault("log.agent", "agent.log")
	viper.SetDefault("log.user_handler", "user_handler.log")
	viper.SetDefault("log.admin_handler", "admin_handler.log")
	viper.SetDefault("log.level", "INFO")
	// db
	viper.SetDefault("db.driver", "sqlite3")
	viper.SetDefault("db.connect", "./data/ignite.db")
	// host
	viper.SetDefault("host.address", "localhost")
	viper.SetDefault("host.from", "5001")
	viper.SetDefault("host.to", "6000")
	// admin
	viper.SetDefault("admin.username", "admin")
	viper.SetDefault("admin.password", "changeme")

	// bind envs
	viper.SetEnvPrefix("ignite")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("viper.Unmarshal error: %v\n", err)
	}

	C.MustCheck()
}

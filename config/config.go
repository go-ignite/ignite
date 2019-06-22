package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/wire"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Set = wire.NewSet(
	Init,
	wire.FieldsOf(new(*Config), "Service", "Server", "Model", "State"),
)

type Server struct {
	Address string `mapstructure:"address"`
}

type Service struct {
	AdminUsername string `mapstructure:"admin_username"`
	AdminPassword string `mapstructure:"admin_password"`
	Secret        string `mapstructure:"secret"`
}

type Model struct {
	Driver  string `mapstructure:"driver"`
	Connect string `mapstructure:"connect"`
	Debug   bool   `mapstructure:"debug"`
}

type State struct {
	SyncInterval        time.Duration `mapstructure:"sync_interval"`
	HeartbeatInterval   time.Duration `mapstructure:"heartbeat_interval"`
	StreamRetryInterval time.Duration `mapstructure:"stream_retry_interval"`
}

type Config struct {
	LogLevel string   `mapstructure:"log_level"`
	Server   *Server  `mapstructure:"server"`
	Service  *Service `mapstructure:"service"`
	Model    *Model   `mapstructure:"model"`
	State    *State   `mapstructure:"state"`
}

func Init() (*Config, error) {
	viper.SetDefault("log_level", "INFO")

	viper.SetDefault("server.address", ":5000")

	viper.SetDefault("service.secret", "ignite")
	viper.SetDefault("service.admin_username", "admin")
	viper.SetDefault("service.admin_password", "changeme")

	viper.SetDefault("model.driver", "sqlite3")
	viper.SetDefault("model.connect", "./data/ignite.db")
	viper.SetDefault("model.debug", false)

	viper.SetDefault("state.heartbeat_interval", time.Second)
	viper.SetDefault("state.sync_interval", 5*time.Second)
	viper.SetDefault("state.stream_retry_interval", 3*time.Second)

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "config: load .env file error")
	}

	// bind envs
	viper.SetEnvPrefix("ignite")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, errors.Wrap(err, "config: unmarshal error")
	}

	if _, err := logrus.ParseLevel(c.LogLevel); err != nil {
		return nil, fmt.Errorf("config: log_level is invalid")
	}

	return c, nil
}

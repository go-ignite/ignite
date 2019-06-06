package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/wire"
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
	JWTSecret     string `mapstructure:"jwt_secret"`
}

type Model struct {
	Driver  string `mapstructure:"driver"`
	Connect string `mapstructure:"connect"`
}

type State struct {
	SyncInterval            time.Duration `mapstructure:"sync_interval"`
	SyncStreamRetryInterval time.Duration `mapstructure:"sync_retry_interval"`
}

type Config struct {
	LogLevel string   `mapstructure:"log_level"`
	Server   *Server  `mapstructure:"server"`
	Service  *Service `mapstructure:"service"`
	Model    *Model   `mapstructure:"model"`
	State    *State   `mapstructure:"state"`
}

func (c *Config) check() error {
	if _, err := logrus.ParseLevel(c.LogLevel); err != nil {
		return fmt.Errorf("config: log_level is invalid")
	}

	return nil
}

func Init() (*Config, error) {
	viper.SetDefault("log_level", "INFO")

	viper.SetDefault("server.address", ":5000")

	viper.SetDefault("service.jwt_secret", "ignite")
	viper.SetDefault("service.admin_username", "admin")
	viper.SetDefault("service.admin_password", "changeme")

	viper.SetDefault("model.driver", "sqlite3")
	viper.SetDefault("model.connect", "./data/ignite.db")

	// bind envs
	viper.SetEnvPrefix("ignite")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, errors.Wrap(err, "config: unmarshal error")
	}

	if err := c.check(); err != nil {
		return nil, err
	}

	return c, nil
}

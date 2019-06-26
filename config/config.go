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
	AdminUsername string        `mapstructure:"admin_username"`
	AdminPassword string        `mapstructure:"admin_password"`
	Secret        string        `mapstructure:"secret"`
	TokenDuration time.Duration `mapstructure:"token_duration"`
}

type Model struct {
	Driver  string `mapstructure:"driver"`
	Connect string `mapstructure:"connect"`
	Debug   bool   `mapstructure:"debug"`
}

type State struct {
	AgentToken          string        `mapstructure:"agent_token"`
	SyncInterval        time.Duration `mapstructure:"sync_interval"`
	HeartbeatInterval   time.Duration `mapstructure:"heartbeat_interval"`
	StreamRetryInterval time.Duration `mapstructure:"stream_retry_interval"`
}

type Config struct {
	LogLevel string  `mapstructure:"log_level"`
	Server   Server  `mapstructure:"server"`
	Service  Service `mapstructure:"service"`
	Model    Model   `mapstructure:"model"`
	State    State   `mapstructure:"state"`
}

var defaultConfig = Config{
	LogLevel: logrus.InfoLevel.String(),
	Server: Server{
		Address: ":5000",
	},
	Service: Service{
		Secret:        "ignite",
		AdminUsername: "admin",
		AdminPassword: "changeme",
		TokenDuration: 24 * time.Hour,
	},
	Model: Model{
		Driver:  "sqlite3",
		Connect: "./data/ignite.db",
		Debug:   false,
	},
	State: State{
		HeartbeatInterval:   time.Second,
		SyncInterval:        5 * time.Second,
		StreamRetryInterval: 3 * time.Second,
		AgentToken:          "ignite-agent",
	},
}

func Init() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "config: load .env file error")
	}

	// bind envs
	viper.SetEnvPrefix("ignite")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	c := defaultConfig
	if err := viper.Unmarshal(&c); err != nil {
		return nil, errors.Wrap(err, "config: unmarshal error")
	}

	lv, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("config: log_level is invalid")
	}
	logrus.SetLevel(lv)

	if c.State.SyncInterval <= 0 {
		return nil, fmt.Errorf("config: state.sync_interval is invalid")
	}

	if c.State.HeartbeatInterval <= 0 {
		return nil, fmt.Errorf("config: state.heartbeat_interval is invalid")
	}

	if c.State.StreamRetryInterval <= 0 {
		return nil, fmt.Errorf("config: state.stream_retry_interval is invalid")
	}

	return &c, nil
}

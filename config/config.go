package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config project config
type Config struct {
	LineCred LineCred
	Server   Server
}

// LineCred for line bot credentials
type LineCred struct {
	Secret string `envconfig:"LINE_CHANNEL_SECRET"`
	Token  string `envconfig:"LINE_CHANNEL_TOKEN"`
}

// Server for server config
type Server struct {
	Addr  string `envconfig:"APP_SERVER_ADDR"`
	Host  string `envconfig:"APP_SERVER_HOST" default:"localhost"`
	Proto string `envconfig:"APP_SERVER_PROTO" default:"http"`
	Port  string `envconfig:"APP_SERVER_PORT" default:"8080"`
	Debug bool   `envconfig:"APP_SERVER_DEBUG" default:"true"`
}

func defualtAddr(c *Config) {
	c.Server.Addr = fmt.Sprintf("%s://%s", c.Server.Proto, c.Server.Host)
}

// Environ load project config
func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	defualtAddr(&cfg)
	return cfg, err
}

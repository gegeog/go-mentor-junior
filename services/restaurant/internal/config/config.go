package config

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Logger   LoggerConfig   `envconfig:"LOGGER"`
	Server   ServerConfig   `envconfig:"HTTP"`
	Postgres PostgresConfig `envconfig:"POSTGRES"`
}

type LoggerConfig struct {
	Level string `envconfig:"LEVEL" default:"DEBUG"`
}

type ServerConfig struct {
	Address         string        `envconfig:"ADDR" default:":8080"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"25s"`
}

type PostgresConfig struct {
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	DB       string `envconfig:"DB"`
	Host     string `envconfig:"HOST"`
}

func (c PostgresConfig) URL() string {
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, "5432"),
		Path:     "/" + c.DB,
		RawQuery: "sslmode=disable",
	}

	return u.String()
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("get config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

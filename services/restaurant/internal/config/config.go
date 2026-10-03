package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Logger LoggerConfig `envconfig:"LOGGER"`
	Server ServerConfig `envconfig:"HTTP"`
}

type LoggerConfig struct {
	Level string `envconfig:"LEVEL" default:"DEBUG"`
}

type ServerConfig struct {
	Address         string        `envconfig:"ADDR" default:":8080"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"25s"`
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

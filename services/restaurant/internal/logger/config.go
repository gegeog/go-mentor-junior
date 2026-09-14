package logger

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	level string `envconfig:"LEVEL" default:"DEBUG"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return cfg
}

func (c Config) GetLevel() slog.Level {
	logLevel := strings.ToLower(c.level)
	switch logLevel {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

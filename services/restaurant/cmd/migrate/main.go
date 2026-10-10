package main

import (
	"os"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/config"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
	"github.com/gegeog/go-mentor-junior/services/restaurant/migrations"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfigMust()
	log := logger.NewLogger(cfg.Logger)

	if len(os.Args) < 2 {
		log.Error("missing command argument")
		os.Exit(1)
	}

	connCfg, err := pgx.ParseConfig(cfg.Postgres.URL())
	if err != nil {
		log.Error("failed to parse postgres connection string", zap.Error(err))
		os.Exit(1)
	}

	if err := migrations.Run(*connCfg, log, os.Args...); err != nil {
		os.Exit(1)
	}
}

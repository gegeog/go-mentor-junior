package migrations

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Run(conn pgx.ConnConfig, logger *zap.Logger, args ...string) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("can not set dialect in goose", zap.Error(err))
		return err
	}

	db := stdlib.OpenDB(conn)
	defer db.Close()

	cmd := args[1]
	ctx := context.Background()
	if err := goose.RunContext(ctx, cmd, db, "migrations", args[2:]...); err != nil {
		logger.Error("can not migrate", zap.Error(err))
		return err
	}

	return nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/handler"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/middleware"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/config"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/provider/postgres"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfigMust()

	logger := logger.NewLogger(cfg.Logger)
	defer logger.Sync()

	ctx := context.Background()

	logger.Info("initializing repo")
	pgxConfig, err := pgxpool.ParseConfig(cfg.Postgres.URL())
	if err != nil {
		logger.Error("failed to parse postgres connection string", zap.Error(err))
		os.Exit(1)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		logger.Error("can not create pgx pool", zap.Error(err))
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("can not ping pool", zap.Error(err))
		os.Exit(1)
	}

	defer pool.Close()

	restaurantsRepository := postgres.NewRestaurantsRepository(pool)

	logger.Info("initializing service")
	restaurantsService := service.NewRestaurantsService(restaurantsRepository)

	logger.Info("initializing transport")
	restaurantsTransport := handler.NewRestaurantHTTPHandler(restaurantsService, logger)

	logger.Info("initializing router")
	mux := api.NewRouter()
	mux.RegisterRoutes(handler.GetKuberStubRoutes(pool)...)
	mux.RegisterRoutes(restaurantsTransport.Routes()...)

	logger.Info("initializing server")
	server := &http.Server{
		Addr: cfg.Server.Address,
		Handler: middleware.ChainMiddleware(
			mux,
			middleware.RequestID(),
			middleware.Trace(logger),
			middleware.Panic(logger),
		),
	}

	//graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		logger.Warn("start HTTP server", zap.String("addr", cfg.Server.Address))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			ch <- fmt.Errorf("HTTP server listen and serve: %w", err)
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			logger.Error("server error", zap.Error(err))
			os.Exit(1)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			cfg.Server.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			closeErr := server.Close()
			logger.Error(
				"failed to shutdown server gracefully",
				zap.Error(errors.Join(err, closeErr)),
			)
		}

		if err := <-ch; err != nil {
			logger.Error(
				"server exited with error",
				zap.Error(err),
			)
		}

		logger.Warn("HTTP server stopped")
	}
}

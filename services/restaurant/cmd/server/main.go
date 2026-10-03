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
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/provider/memory"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfigMust()

	logger := logger.NewLogger(cfg.Logger)
	defer logger.Sync()

	logger.Info("initializing restaurants service")
	restaurantsRepository := memory.NewRestaurantsRepository()
	restaurantsService := service.NewRestaurantsService(restaurantsRepository)
	restaurantsTransport := handler.NewRestaurantHTTPHandler(restaurantsService, logger)

	logger.Info("initializing router")
	mux := api.NewRouter()
	mux.RegisterRoutes(handler.GetKuberStubRoutes()...)
	mux.RegisterRoutes(restaurantsTransport.Routes()...)

	logger.Info("initializing server")
	server := &http.Server{
		Addr: cfg.Server.Address,
		Handler: middleware.ChainMiddleware(
			mux,
			middleware.RequestID(),
			// TODO: я убрал логгер из контекста и передаю теперь прямо в мидлварю и в хэндлеры...
			// и мне кажется это решение более громоздким...правки хэндлеров??
			// хотя наверняка позже пойму преимущества когда будем с otel работать
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

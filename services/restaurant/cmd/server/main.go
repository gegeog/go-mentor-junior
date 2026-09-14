package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/handler"
	restaurants_transport "github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/handler/restaurants"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/middleware"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/app"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
	restaurants_repository "github.com/gegeog/go-mentor-junior/services/restaurant/internal/provider/memory/restaurants"
	restaurants_service "github.com/gegeog/go-mentor-junior/services/restaurant/internal/service/restaurants"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger := logger.NewLogger(logger.NewConfigMust())

	logger.Debug("initializing server...")
	server := app.NewHTTPServer(
		app.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Panic(),
	)

	logger.Debug("initializing restaurants service")
	restaurantsRepository := restaurants_repository.NewRestaurantsRepository()
	restaurantsService := restaurants_service.NewRestaurantsService(restaurantsRepository)
	restaurantsTransport := restaurants_transport.NewRestaurantHTTPHandler(restaurantsService)

	logger.Debug("initializing router")
	mux := app.NewRouter()
	mux.RegisterRoutes(handler.GetKuberStubRoutes()...)
	mux.RegisterRoutes(restaurantsTransport.Routes()...)

	if err := server.Run(ctx, mux); err != nil {
		logger.Error("HTTP server run error", slog.Any("error", err))
	}
}

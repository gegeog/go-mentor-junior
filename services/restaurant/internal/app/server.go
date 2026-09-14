package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/middleware"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

type HTTPServer struct {
	config Config
	log    *logger.Logger

	middleware []middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *logger.Logger,
	middleware ...middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) Run(
	ctx context.Context,
	router *Router,
) error {
	mux := middleware.ChainMiddleware(router, s.middleware...)
	server := &http.Server{
		Addr:    s.config.Address,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", slog.String("addr", s.config.Address))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("HTTP server listen and serve: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}

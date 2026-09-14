package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			log.With(
				slog.String("request_id", requestID),
				slog.String("url", r.URL.String()),
			)

			ctx := logger.ToContext(
				r.Context(),
				log,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

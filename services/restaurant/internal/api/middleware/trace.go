package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseWriter := response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				"incoming HTTP request",
				slog.String("http_method", r.Method),
				slog.Time("time", before),
			)

			next.ServeHTTP(responseWriter, r)

			log.Debug(
				"done HTTP request",
				slog.Int("status_code", responseWriter.GetStatusCode()),
				slog.Duration("latency", time.Since(before)),
			)
		})
	}
}

package middleware

import (
	"net/http"
	"time"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"go.uber.org/zap"
)

func Trace(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseWriter := response.NewResponseWriter(w)

			before := time.Now()
			logger.Debug(
				"incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before),
			)

			next.ServeHTTP(responseWriter, r)

			logger.Debug(
				"done HTTP request",
				zap.Int("status_code", responseWriter.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

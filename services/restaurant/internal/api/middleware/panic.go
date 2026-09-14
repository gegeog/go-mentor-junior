package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)

			defer func() {
				if p := recover(); p != nil {
					log.Error(
						"during handle HTTP request got panic",
						slog.Any("error", p),
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

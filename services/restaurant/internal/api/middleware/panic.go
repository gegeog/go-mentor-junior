package middleware

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"go.uber.org/zap"
)

func Panic(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					response.PanicResponse(
						logger,
						w,
						p,
						"during handle HTTP request got panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func RestaurantID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := response.NewHTTPResponseHandler(log, w)

			restaurantID := r.Header.Get("X-Restaurant-ID")
			if restaurantID == "" {
				responseHandler.ErrorResponse(
					domain.ErrInvalidRestaurant,
					"failed to get X-Restaurant-ID from request header",
				)

				return
			}

			if restaurantID != r.PathValue("restaurant_id") {
				responseHandler.ErrorResponse(
					domain.ErrInvalidRestaurant,
					"X-Restaurant-ID provided in headers not match with path value",
				)

				return
			}

			next.ServeHTTP(w, r)

		})
	}
}

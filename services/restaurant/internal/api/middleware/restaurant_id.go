package middleware

import (
	"fmt"
	"net/http"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"go.uber.org/zap"
)

func RestaurantID(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerRaw := r.Header.Get(request.HeaderRestaurantID)
			if headerRaw == "" {
				response.ErrorResponse(
					logger,
					w,
					domain.ErrInvalidRestaurant,
					fmt.Sprintf(
						"failed to get %s from request header",
						request.HeaderRestaurantID,
					),
				)

				return
			}

			headerID, err := uuid.Parse(headerRaw)
			if err != nil || headerID == uuid.Nil() {
				response.ErrorResponse(
					logger,
					w,
					domain.ErrInvalidRestaurant,
					fmt.Sprintf("invalid %s", request.HeaderRestaurantID),
				)

				return
			}

			pathID, err := request.ParseRestaurantID(r)
			if err != nil {
				response.ErrorResponse(logger, w, err, "invalid restaurant_id")
				return
			}

			if headerID != pathID {
				response.ErrorResponse(
					logger,
					w,
					domain.ErrRestaurantAccessDenied,
					fmt.Sprintf(
						"%s provided in headers not match with path value",
						request.HeaderRestaurantID,
					),
				)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

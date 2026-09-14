package restaurants_transport

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func (h *RestaurantHTTPHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id := r.PathValue("restaurant_id")
	if id == "" {
		responseHandler.ErrorResponse(
			domain.ErrInvalidRestaurant,
			"failed to get restaurantID path value",
		)

		return
	}

	restaurantDomain, err := h.restaurantService.GetRestaurant(id, ctx)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get restaurant",
		)

		return
	}
	response := restaurantDTOFromDomain(restaurantDomain)
	responseHandler.JSONResponse(response, http.StatusOK)
}

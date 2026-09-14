package restaurants_transport

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func (h *RestaurantHTTPHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	//TODO: add pagination(limit + offset)
	userDomains, err := h.restaurantService.GetRestaurants(ctx)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get restaurants",
		)

		return
	}

	restaurantsDTO := make([]RestaurantResponseDTO, len(userDomains))
	for i, userDomain := range userDomains {
		restaurantsDTO[i] = restaurantDTOFromDomain(userDomain)
	}

	var restaurantsResponse = map[string][]RestaurantResponseDTO{
		"restaurants": restaurantsDTO,
	}

	responseHandler.JSONResponse(restaurantsResponse, http.StatusOK)
}

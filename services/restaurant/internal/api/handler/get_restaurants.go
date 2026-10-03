package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
)

func (h *RestaurantHTTPHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	log := h.requestLogger(r)

	//TODO: add pagination(limit + offset)
	restaurantsDomains, err := h.restaurantService.GetRestaurants(r.Context())
	if err != nil {
		response.ErrorResponse(log, w, err, "failed to get restaurants")
		return
	}

	restaurantsDTO := make([]RestaurantResponseDTO, len(restaurantsDomains))
	for i, userDomain := range restaurantsDomains {
		restaurantsDTO[i] = restaurantDTOFromDomain(userDomain)
	}

	restaurantsResponse := RestaurantsResponseDTO{
		Restaurants: restaurantsDTO,
	}

	response.JSONResponse(log, w, restaurantsResponse, http.StatusOK)
}

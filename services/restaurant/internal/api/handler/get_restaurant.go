package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
)

func (h *RestaurantHTTPHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	log := h.requestLogger(r)

	id, err := request.ParseRestaurantID(r)
	if err != nil {
		response.ErrorResponse(log, w, err, "invalid restaurant_id")
		return
	}

	restaurantDomain, err := h.restaurantService.GetRestaurant(r.Context(), id)
	if err != nil {
		response.ErrorResponse(log, w, err, "failed to get restaurant")
		return
	}
	rsp := restaurantDTOFromDomain(restaurantDomain)
	response.JSONResponse(log, w, rsp, http.StatusOK)
}

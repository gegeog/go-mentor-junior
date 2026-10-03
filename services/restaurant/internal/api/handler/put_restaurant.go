package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (h *RestaurantHTTPHandler) PutRestaurant(w http.ResponseWriter, r *http.Request) {
	log := h.requestLogger(r)

	restaurantID, err := request.ParseRestaurantID(r)
	if err != nil {
		response.ErrorResponse(log, w, err, "invalid restaurant_id")
		return
	}

	req, err := request.DecodeAndValidate[PutRestaurantRequest](r)
	if err != nil {
		response.ErrorResponse(log, w, domain.ErrInvalidRestaurant, "failed to decode and validate request")
		return
	}

	restaurantDomain := domain.NewRestaurant(
		restaurantID,
		req.Name,
		*req.AcceptingOrders,
		req.MinimumOrderMinor,
		req.Currency,
	)

	restaurantDomain, err = h.restaurantService.PutRestaurant(r.Context(), restaurantDomain)
	if err != nil {
		response.ErrorResponse(log, w, err, "failed to put restaurant")
		return
	}

	rsp := restaurantDTOFromDomain(restaurantDomain)
	response.JSONResponse(log, w, rsp, http.StatusOK)
}

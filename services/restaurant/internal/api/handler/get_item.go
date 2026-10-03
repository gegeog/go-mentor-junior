package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
)

func (h *RestaurantHTTPHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	log := h.requestLogger(r)

	restaurantID, err := request.ParseRestaurantID(r)
	if err != nil {
		response.ErrorResponse(log, w, err, "invalid restaurant_id")
		return
	}

	itemID, err := request.ParseMenuItemID(r)
	if err != nil {
		response.ErrorResponse(log, w, err, "invalid menu_item_id")
		return
	}

	restaurant, err := h.restaurantService.GetItem(
		r.Context(),
		restaurantID,
		itemID,
	)
	if err != nil {
		response.ErrorResponse(log, w, err, "failed to get menu item")
		return
	}

	rsp := menuItemDTOFromDomain(restaurant)
	response.JSONResponse(log, w, rsp, http.StatusOK)
}

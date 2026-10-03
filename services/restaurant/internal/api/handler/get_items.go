package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
)

func (h *RestaurantHTTPHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	log := h.requestLogger(r)

	restaurantID, err := request.ParseRestaurantID(r)
	if err != nil {
		response.ErrorResponse(log, w, err, "invalid restaurant_id")
		return
	}

	items, err := h.restaurantService.GetItems(r.Context(), restaurantID)
	if err != nil {
		response.ErrorResponse(log, w, err, "failed to get menu items")
		return
	}

	itemsDTO := make([]MenuItemResponseDTO, len(items))
	for i, v := range items {
		itemsDTO[i] = menuItemDTOFromDomain(v)
	}

	rsp := MenuItemsResponseDTO{
		RestaurantID: restaurantID.String(),
		Items:        itemsDTO,
	}

	response.JSONResponse(log, w, rsp, http.StatusOK)
}

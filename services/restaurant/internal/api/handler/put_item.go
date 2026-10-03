package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (h *RestaurantHTTPHandler) PutMenuItem(w http.ResponseWriter, r *http.Request) {
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

	req, err := request.DecodeAndValidate[PutMenuItemRequest](r)
	if err != nil {
		response.ErrorResponse(log, w, domain.ErrInvalidMenuItem, "failed to decode and validate request")
		return
	}

	itemDomain := domain.NewMenuItem(
		itemID,
		req.Name,
		req.Description,
		req.PriceMinor,
		req.Currency,
		*req.Available,
	)

	itemDomain, err = h.restaurantService.PutMenuItem(
		r.Context(),
		restaurantID,
		itemDomain,
	)
	if err != nil {
		response.ErrorResponse(log, w, err, "failed to put menu item")
		return
	}

	rsp := menuItemDTOFromDomain(itemDomain)
	response.JSONResponse(log, w, rsp, http.StatusOK)
}

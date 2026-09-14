package restaurants_transport

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

type PutMenuItemRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	PriceMinor  int64  `json:"price_minor" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
	Available   bool   `json:"available" validate:"required"`
}

type MenuItemResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceMinor  int64  `json:"price_minor"`
	Currency    string `json:"currency"`
	Available   bool   `json:"available"`
}

func (h *RestaurantHTTPHandler) PutMenuItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	restaurantID := r.PathValue("restaurant_id")
	if restaurantID == "" {
		responseHandler.ErrorResponse(
			domain.ErrInvalidRestaurant,
			"failed to get restaurantID path value",
		)

		return
	}

	itemID := r.PathValue("menu_item_id")
	if itemID == "" {
		responseHandler.ErrorResponse(
			domain.ErrInvalidMenuItem,
			"failed to get menu_item_id path value",
		)

		return
	}

	var req PutMenuItemRequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		responseHandler.ErrorResponse(
			domain.ErrInvalidMenuItem,
			"failed to decode and validate request",
		)

		return
	}

	itemDomain := domain.NewMenuItem(
		itemID,
		req.Name,
		req.Description,
		req.PriceMinor,
		req.Currency,
		req.Available,
	)

	itemDomain, err := h.restaurantService.PutMenuItem(ctx, restaurantID, itemDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to put menu item",
		)
		return
	}

	response := menuItemDTOFromDomain(itemDomain)
	responseHandler.JSONResponse(response, http.StatusOK)
}

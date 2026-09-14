package restaurants_transport

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

func (h *RestaurantHTTPHandler) GetItem(w http.ResponseWriter, r *http.Request) {
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

	restaurant, err := h.restaurantService.GetItem(ctx, restaurantID, itemID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get menu item",
		)

		return
	}

	response := menuItemDTOFromDomain(restaurant)
	responseHandler.JSONResponse(response, http.StatusOK)
}

package restaurants_transport

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

type MenuItemsResponseDTO struct {
	RestaurantID string                `json:"restaurant_id"`
	Items        []MenuItemResponseDTO `json:"items"`
}

func (h *RestaurantHTTPHandler) GetItems(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.restaurantService.GetItems(ctx, restaurantID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get menu items",
		)

		return
	}

	itemsDTO := make([]MenuItemResponseDTO, len(items))
	for i, v := range items {
		itemsDTO[i] = menuItemDTOFromDomain(v)
	}

	rsp := MenuItemsResponseDTO{
		RestaurantID: restaurantID,
		Items:        itemsDTO,
	}

	responseHandler.JSONResponse(rsp, http.StatusOK)
}

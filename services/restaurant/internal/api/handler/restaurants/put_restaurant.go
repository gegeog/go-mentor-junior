package restaurants_transport

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/request"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

type PutRestaurantRequest struct {
	Name               string `json:"name" validate:"required"`
	AcceptingOrders    bool   `json:"accepting_orders" validate:"required"`
	MinimumOrderAmount int64  `json:"minimum_order_amount" validate:"required"`
	Currency           string `json:"currency" validate:"required"`
}

type RestaurantResponseDTO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	AcceptingOrders    bool   `json:"accepting_orders"`
	MinimumOrderAmount int64  `json:"minimum_order_amount"`
	Currency           string `json:"currency"`
}

func (h *RestaurantHTTPHandler) PutRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id := r.PathValue("restaurant_id")
	if id == "" {
		responseHandler.ErrorResponse(
			domain.ErrInvalidRestaurant,
			"failed to get restaurantID path value",
		)

		return
	}

	var req PutRestaurantRequest
	if err := request.DecodeAndValidate(r, &req); err != nil {
		responseHandler.ErrorResponse(
			domain.ErrInvalidRestaurant,
			"failed to decode and validate request",
		)

		return
	}

	restaurantDomain := domain.NewRestaurant(
		id,
		req.Name,
		req.AcceptingOrders,
		req.MinimumOrderAmount,
		req.Currency,
	)

	restaurantDomain, err := h.restaurantService.PutRestaurant(ctx, restaurantDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to put restaurant",
		)

		return
	}

	response := restaurantDTOFromDomain(restaurantDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}

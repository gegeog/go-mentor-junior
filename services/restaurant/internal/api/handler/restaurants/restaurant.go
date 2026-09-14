package restaurants_transport

import (
	"context"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/middleware"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/app"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

type RestaurantHTTPHandler struct {
	restaurantService RestaurantService
}

type RestaurantService interface {
	PutRestaurant(context.Context, domain.Restaurant) (domain.Restaurant, error)
	GetRestaurants(context.Context) ([]domain.Restaurant, error)
	GetRestaurant(string, context.Context) (domain.Restaurant, error)
	PutMenuItem(context.Context, string, domain.MenuItem) (domain.MenuItem, error)
	GetItems(context.Context, string) ([]domain.MenuItem, error)
	GetItem(context.Context, string, string) (domain.MenuItem, error)
}

func NewRestaurantHTTPHandler(
	restaurantService RestaurantService,
) *RestaurantHTTPHandler {
	return &RestaurantHTTPHandler{
		restaurantService: restaurantService,
	}
}

func (h *RestaurantHTTPHandler) Routes() []app.Route {
	return []app.Route{
		{
			// TODO: чет хуй пойму какое должно быть поведение при изменении currency ресторана. например изначально создали рестик с currency "RUB", потом накидали туда блюд с той же валютой. а затем решили заменить currency рестика на "USD", то надо же и все блюда менять на "USD"? но прайс же по-прежнему выставлен в "RUB", то есть надо ещё куда-то ходить и брать актуальный курс валют?
			Path:    "/restaurants/{restaurant_id}",
			Method:  http.MethodPut,
			Handler: h.PutRestaurant,
			Middleware: []middleware.Middleware{
				middleware.RestaurantID(),
			},
		},
		{
			Path:    "/restaurants",
			Method:  http.MethodGet,
			Handler: h.GetRestaurants,
		},
		{
			Path:    "/restaurants/{restaurant_id}",
			Method:  http.MethodGet,
			Handler: h.GetRestaurant,
		},
		{
			Path:    "/restaurants/{restaurant_id}/menu/{menu_item_id}",
			Method:  http.MethodPut,
			Handler: h.PutMenuItem,
			Middleware: []middleware.Middleware{
				middleware.RestaurantID(),
			},
		},
		{
			Path:    "/restaurants/{restaurant_id}/menu",
			Method:  http.MethodGet,
			Handler: h.GetItems,
		},
		{
			Path:    "/restaurants/{restaurant_id}/menu/{menu_item_id}",
			Method:  http.MethodGet,
			Handler: h.GetItem,
		},
	}
}

func restaurantDTOFromDomain(restaurantDomain domain.Restaurant) RestaurantResponseDTO {
	return RestaurantResponseDTO{
		ID:                 restaurantDomain.ID,
		Name:               restaurantDomain.Name,
		AcceptingOrders:    restaurantDomain.AcceptingOrders,
		MinimumOrderAmount: restaurantDomain.MinimumOrderMminor,
		Currency:           restaurantDomain.Currency,
	}
}

func menuItemDTOFromDomain(itemDomain domain.MenuItem) MenuItemResponseDTO {
	return MenuItemResponseDTO{
		ID:          itemDomain.ID,
		Name:        itemDomain.Name,
		Description: itemDomain.Description,
		PriceMinor:  itemDomain.PriceMinor,
		Currency:    itemDomain.Currency,
		Available:   itemDomain.Available,
	}
}

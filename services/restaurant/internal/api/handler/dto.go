package handler

import "github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"

type MenuItemsResponseDTO struct {
	RestaurantID string                `json:"restaurant_id"`
	Items        []MenuItemResponseDTO `json:"items"`
}

type RestaurantsResponseDTO struct {
	Restaurants []RestaurantResponseDTO `json:"restaurants"`
}

type PutMenuItemRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description" validate:"required"`
	PriceMinor  int64  `json:"price_minor" validate:"required,gt=0"`
	Currency    string `json:"currency" validate:"required,min=1"`
	Available   *bool  `json:"available" validate:"required"`
}

type MenuItemResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceMinor  int64  `json:"price_minor"`
	Currency    string `json:"currency"`
	Available   bool   `json:"available"`
}

type PutRestaurantRequest struct {
	Name              string `json:"name" validate:"required,min=1"`
	AcceptingOrders   *bool  `json:"accepting_orders" validate:"required"`
	MinimumOrderMinor int64  `json:"minimum_order_minor" validate:"required,gt=0"`
	Currency          string `json:"currency" validate:"required,min=1"`
}

type RestaurantResponseDTO struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	AcceptingOrders   bool   `json:"accepting_orders"`
	MinimumOrderMinor int64  `json:"minimum_order_minor"`
	Currency          string `json:"currency"`
}

func restaurantDTOFromDomain(restaurantDomain domain.Restaurant) RestaurantResponseDTO {
	return RestaurantResponseDTO{
		ID:                restaurantDomain.ID.String(),
		Name:              restaurantDomain.Name,
		AcceptingOrders:   restaurantDomain.AcceptingOrders,
		MinimumOrderMinor: restaurantDomain.MinimumOrderMinor,
		Currency:          restaurantDomain.Currency,
	}
}

func menuItemDTOFromDomain(itemDomain domain.MenuItem) MenuItemResponseDTO {
	return MenuItemResponseDTO{
		ID:          itemDomain.ID.String(),
		Name:        itemDomain.Name,
		Description: itemDomain.Description,
		PriceMinor:  itemDomain.PriceMinor,
		Currency:    itemDomain.Currency,
		Available:   itemDomain.Available,
	}
}

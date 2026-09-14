package restaurants_repository

import "github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"

type RestaurantModel struct {
	Name              string
	AcceptingOrders   bool
	MinimumOrderMinor int64
	Currency          string
}

type MenuItemModel struct {
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	Available   bool
}

func restaurantModelFromDomain(restaurantDomain domain.Restaurant) RestaurantModel {
	return RestaurantModel{
		Name:              restaurantDomain.Name,
		AcceptingOrders:   restaurantDomain.AcceptingOrders,
		MinimumOrderMinor: restaurantDomain.MinimumOrderMminor,
		Currency:          restaurantDomain.Currency,
	}
}

func menuItemModelFromDomain(itemDomain domain.MenuItem) MenuItemModel {
	return MenuItemModel{
		Name:        itemDomain.Name,
		Description: itemDomain.Description,
		PriceMinor:  itemDomain.PriceMinor,
		Currency:    itemDomain.Currency,
		Available:   itemDomain.Available,
	}
}

func restaurantDomainFromModel(
	id string,
	model RestaurantModel,
) domain.Restaurant {
	return domain.NewRestaurant(
		id,
		model.Name,
		model.AcceptingOrders,
		model.MinimumOrderMinor,
		model.Currency,
	)
}

func itemDomainFromModel(
	id string,
	model MenuItemModel,
) domain.MenuItem {
	return domain.NewMenuItem(
		id,
		model.Name,
		model.Description,
		model.PriceMinor,
		model.Currency,
		model.Available,
	)
}

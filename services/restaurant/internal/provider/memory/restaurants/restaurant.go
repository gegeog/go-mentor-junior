package restaurants_repository

import (
	"sync"
)

type RestaurantsRepository struct {
	restaurants    map[string]RestaurantModel
	items          map[string]map[string]MenuItemModel
	restaurantsMtx sync.RWMutex
	itemsMtx       sync.RWMutex
}

func NewRestaurantsRepository() *RestaurantsRepository {
	return &RestaurantsRepository{
		restaurants: make(map[string]RestaurantModel),
		items:       make(map[string]map[string]MenuItemModel),
	}
}

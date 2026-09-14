package restaurants_repository

import (
	"context"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) CreateOrReplaceRestaurant(
	ctx context.Context,
	restaurant domain.Restaurant,
) (domain.Restaurant, error) {
	r.restaurantsMtx.Lock()
	defer r.restaurantsMtx.Unlock()

	r.restaurants[restaurant.ID] = restaurantModelFromDomain(restaurant)

	r.itemsMtx.Lock()
	defer r.itemsMtx.Unlock()
	r.items[restaurant.ID] = make(map[string]MenuItemModel)

	return restaurant, nil
}

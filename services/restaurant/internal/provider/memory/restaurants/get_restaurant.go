package restaurants_repository

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) GetRestaurant(
	id string,
	ctx context.Context,
) (domain.Restaurant, error) {
	r.restaurantsMtx.RLock()
	defer r.restaurantsMtx.RUnlock()

	restaurant, ok := r.restaurants[id]
	if !ok {
		return domain.Restaurant{}, fmt.Errorf(
			"restaurant with id=%s: %w",
			id,
			domain.ErrRestaurantNotFound,
		)
	}

	return restaurantDomainFromModel(id, restaurant), nil
}

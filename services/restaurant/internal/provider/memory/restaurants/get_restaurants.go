package restaurants_repository

import (
	"context"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) GetRestaurants(
	ctx context.Context,
) ([]domain.Restaurant, error) {
	r.restaurantsMtx.RLock()
	defer r.restaurantsMtx.RUnlock()

	restaurants := make([]domain.Restaurant, 0, len(r.restaurants))
	for k, v := range r.restaurants {
		restaurantDomain := restaurantDomainFromModel(k, v)
		restaurants = append(restaurants, restaurantDomain)
	}

	return restaurants, nil
}

package restaurants_repository

import (
	"context"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) PutMenuItem(
	ctx context.Context,
	restaurantID string,
	itemDomain domain.MenuItem,
) (domain.MenuItem, error) {
	_, err := r.GetRestaurant(restaurantID, ctx)
	if err != nil {
		return domain.MenuItem{}, err
	}

	r.itemsMtx.Lock()
	defer r.itemsMtx.Unlock()

	r.items[restaurantID][itemDomain.ID] = menuItemModelFromDomain(itemDomain)

	return itemDomain, nil
}

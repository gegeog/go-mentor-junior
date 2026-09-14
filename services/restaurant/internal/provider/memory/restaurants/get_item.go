package restaurants_repository

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) GetMenuItem(
	restaurantID string,
	itemID string,
	ctx context.Context,
) (domain.MenuItem, error) {
	_, err := r.GetRestaurant(restaurantID, ctx)
	if err != nil {
		return domain.MenuItem{}, err
	}

	r.itemsMtx.RLock()
	defer r.itemsMtx.RUnlock()

	item, ok := r.items[restaurantID][itemID]
	if !ok {
		return domain.MenuItem{}, fmt.Errorf("not found: %w", domain.ErrMenuItemNotFound)
	}

	return itemDomainFromModel(itemID, item), nil
}

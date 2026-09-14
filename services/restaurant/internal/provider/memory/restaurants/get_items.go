package restaurants_repository

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) GetItems(
	ctx context.Context,
	restaurantID string,
) ([]domain.MenuItem, error) {
	_, err := r.GetRestaurant(restaurantID, ctx)
	if err != nil {
		return []domain.MenuItem{}, nil
	}

	r.itemsMtx.RLock()
	defer r.itemsMtx.RUnlock()

	menuItems := make([]domain.MenuItem, 0, len(r.items[restaurantID]))
	for k, v := range r.items[restaurantID] {
		menuItems = append(menuItems, itemDomainFromModel(k, v))
	}

	fmt.Println(menuItems)

	return menuItems, nil
}

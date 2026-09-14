package restaurants_service

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) GetItem(
	ctx context.Context,
	restaurantID string,
	itemID string) (domain.MenuItem, error) {
	item, err := s.restaurantsRepository.GetMenuItem(restaurantID, itemID, ctx)
	if err != nil {
		return domain.MenuItem{}, fmt.Errorf(
			"get menu item from repository: %w",
			err,
		)
	}

	return item, nil
}

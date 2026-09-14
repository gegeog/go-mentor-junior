package restaurants_service

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) GetItems(
	ctx context.Context,
	restaurantID string,
) ([]domain.MenuItem, error) {
	items, err := s.restaurantsRepository.GetItems(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("get menu items: %w", err)
	}

	return items, nil
}

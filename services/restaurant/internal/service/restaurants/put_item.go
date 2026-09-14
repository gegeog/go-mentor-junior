package restaurants_service

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) PutMenuItem(
	ctx context.Context,
	restaurantID string,
	menuItem domain.MenuItem,
) (domain.MenuItem, error) {
	if err := menuItem.Validate(); err != nil {
		return domain.MenuItem{}, fmt.Errorf("validate item domain: %w", err)
	}

	restaurant, err := s.restaurantsRepository.GetRestaurant(restaurantID, ctx)
	if err != nil {
		return domain.MenuItem{}, fmt.Errorf("get restaurant: %w", err)
	}

	if restaurant.Currency != menuItem.Currency {
		return domain.MenuItem{}, fmt.Errorf("currency mismatch: %w", domain.ErrCurrencyMismatch)
	}

	menuItem, err = s.restaurantsRepository.PutMenuItem(ctx, restaurantID, menuItem)
	if err != nil {
		return domain.MenuItem{}, fmt.Errorf("put menu item: %w", err)
	}

	return menuItem, nil
}

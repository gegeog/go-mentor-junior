package service

import (
	"context"
	"fmt"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) GetItem(
	ctx context.Context,
	restaurantID uuid.UUID,
	itemID uuid.UUID,
) (domain.MenuItem, error) {
	item, err := s.restaurantsRepository.GetMenuItem(ctx, restaurantID, itemID)
	if err != nil {
		return domain.MenuItem{}, fmt.Errorf(
			"get menu item from repository: %w",
			err,
		)
	}

	return item, nil
}

func (s *RestaurantsService) GetItems(
	ctx context.Context,
	restaurantID uuid.UUID,
) ([]domain.MenuItem, error) {
	items, err := s.restaurantsRepository.GetMenuItems(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("get menu items: %w", err)
	}

	return items, nil
}

func (s *RestaurantsService) PutMenuItem(
	ctx context.Context,
	restaurantID uuid.UUID,
	menuItem domain.MenuItem,
) (domain.MenuItem, error) {
	restaurant, err := s.restaurantsRepository.GetRestaurant(ctx, restaurantID)
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

package restaurants_service

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) PutRestaurant(
	ctx context.Context,
	restaurant domain.Restaurant,
) (domain.Restaurant, error) {
	if err := restaurant.Validate(); err != nil {
		return domain.Restaurant{}, fmt.Errorf("validate restaurant domain: %w", err)
	}

	restaurant, err := s.restaurantsRepository.CreateOrReplaceRestaurant(ctx, restaurant)
	if err != nil {
		return domain.Restaurant{}, fmt.Errorf("create or replace restaurant: %w", err)
	}

	return restaurant, nil
}

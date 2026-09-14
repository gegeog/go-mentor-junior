package restaurants_service

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) GetRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	restaurants, err := s.restaurantsRepository.GetRestaurants(ctx)
	if err != nil {
		return nil, fmt.Errorf("get restaurants: %w", err)
	}

	return restaurants, nil
}

package restaurants_service

import (
	"context"
	"fmt"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (s *RestaurantsService) GetRestaurant(
	id string,
	ctx context.Context,
) (domain.Restaurant, error) {
	restaurant, err := s.restaurantsRepository.GetRestaurant(id, ctx)
	if err != nil {
		return domain.Restaurant{}, fmt.Errorf(
			"get restaurant from repository: %w",
			err,
		)
	}

	return restaurant, nil
}

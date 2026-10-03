package service

import (
	"context"
	"fmt"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

type RestaurantsService struct {
	restaurantsRepository RestaurantsRepository
}

//go:generate go tool minimock -i RestaurantsRepository -o restaurants_repository_mock_test.go
type RestaurantsRepository interface {
	CreateOrReplaceRestaurant(
		ctx context.Context,
		restaurant domain.Restaurant,
	) (domain.Restaurant, error)

	GetRestaurants(
		ctx context.Context,
	) ([]domain.Restaurant, error)

	GetRestaurant(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Restaurant, error)

	PutMenuItem(
		ctx context.Context,
		restaurantID uuid.UUID,
		itemDomain domain.MenuItem,
	) (domain.MenuItem, error)

	GetMenuItems(
		ctx context.Context,
		restaurantID uuid.UUID,
	) ([]domain.MenuItem, error)

	GetMenuItem(
		ctx context.Context,
		restaurantID uuid.UUID,
		itemID uuid.UUID,
	) (domain.MenuItem, error)
}

func NewRestaurantsService(
	restaurantsRepository RestaurantsRepository,
) *RestaurantsService {
	return &RestaurantsService{
		restaurantsRepository: restaurantsRepository,
	}
}

func (s *RestaurantsService) GetRestaurant(
	ctx context.Context,
	id uuid.UUID,
) (domain.Restaurant, error) {
	restaurant, err := s.restaurantsRepository.GetRestaurant(ctx, id)
	if err != nil {
		return domain.Restaurant{}, fmt.Errorf(
			"get restaurant from repository: %w",
			err,
		)
	}

	return restaurant, nil
}

func (s *RestaurantsService) GetRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	restaurants, err := s.restaurantsRepository.GetRestaurants(ctx)
	if err != nil {
		return nil, fmt.Errorf("get restaurants: %w", err)
	}

	return restaurants, nil
}

func (s *RestaurantsService) PutRestaurant(
	ctx context.Context,
	restaurant domain.Restaurant,
) (domain.Restaurant, error) {
	if r, err := s.restaurantsRepository.GetRestaurant(ctx, restaurant.ID); err == nil {
		if r.Currency != restaurant.Currency {
			return domain.Restaurant{}, fmt.Errorf("create or replace restaurant: %w", domain.ErrCurrencyMismatch)
		}
	}

	restaurant, err := s.restaurantsRepository.CreateOrReplaceRestaurant(ctx, restaurant)
	if err != nil {
		return domain.Restaurant{}, fmt.Errorf("create or replace restaurant: %w", err)
	}

	return restaurant, nil
}

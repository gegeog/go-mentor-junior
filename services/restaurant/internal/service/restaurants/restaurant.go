package restaurants_service

import (
	"context"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

type RestaurantsService struct {
	restaurantsRepository RestaurantsRepository
}

type RestaurantsRepository interface {
	CreateOrReplaceRestaurant(
		ctx context.Context,
		restaurant domain.Restaurant,
	) (domain.Restaurant, error)

	GetRestaurants(
		ctx context.Context,
	) ([]domain.Restaurant, error)

	GetRestaurant(
		id string,
		ctx context.Context,
	) (domain.Restaurant, error)

	PutMenuItem(
		context.Context,
		string,
		domain.MenuItem,
	) (domain.MenuItem, error)

	GetItems(
		context.Context,
		string,
	) ([]domain.MenuItem, error)

	GetMenuItem(
		restaurantID string,
		itemID string,
		ctx context.Context,
	) (domain.MenuItem, error)
}

func NewRestaurantsService(
	restaurantsRepository RestaurantsRepository,
) *RestaurantsService {
	return &RestaurantsService{
		restaurantsRepository: restaurantsRepository,
	}
}

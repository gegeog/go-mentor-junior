package memory

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

type RestaurantsRepository struct {
	restaurants    map[uuid.UUID]RestaurantModel
	items          map[uuid.UUID]map[uuid.UUID]MenuItemModel
	restaurantsMtx sync.RWMutex
	itemsMtx       sync.RWMutex
}

func NewRestaurantsRepository() *RestaurantsRepository {
	return &RestaurantsRepository{
		restaurants: make(map[uuid.UUID]RestaurantModel),
		items:       make(map[uuid.UUID]map[uuid.UUID]MenuItemModel),
	}
}

func (r *RestaurantsRepository) GetRestaurant(
	ctx context.Context,
	id uuid.UUID,
) (domain.Restaurant, error) {
	r.restaurantsMtx.RLock()
	defer r.restaurantsMtx.RUnlock()

	restaurant, ok := r.restaurants[id]
	if !ok {
		return domain.Restaurant{}, fmt.Errorf(
			"restaurant with id=%s: %w",
			id,
			domain.ErrRestaurantNotFound,
		)
	}

	return restaurantDomainFromModel(id, restaurant), nil
}

func (r *RestaurantsRepository) GetRestaurants(
	ctx context.Context,
) ([]domain.Restaurant, error) {
	r.restaurantsMtx.RLock()
	defer r.restaurantsMtx.RUnlock()

	keys := make([]uuid.UUID, 0, len(r.restaurants))
	for k := range r.restaurants {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(a, b uuid.UUID) int {
		return strings.Compare(a.String(), b.String())
	})

	restaurants := make([]domain.Restaurant, 0, len(keys))
	for _, k := range keys {
		restaurants = append(restaurants, restaurantDomainFromModel(k, r.restaurants[k]))
	}

	return restaurants, nil
}

func (r *RestaurantsRepository) CreateOrReplaceRestaurant(
	ctx context.Context,
	restaurant domain.Restaurant,
) (domain.Restaurant, error) {
	r.restaurantsMtx.Lock()
	defer r.restaurantsMtx.Unlock()

	r.restaurants[restaurant.ID] = restaurantModelFromDomain(restaurant)

	r.itemsMtx.Lock()
	defer r.itemsMtx.Unlock()

	if _, ok := r.items[restaurant.ID]; !ok {
		r.items[restaurant.ID] = make(map[uuid.UUID]MenuItemModel)
	}

	return restaurant, nil
}

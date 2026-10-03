package memory

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

func (r *RestaurantsRepository) GetMenuItem(
	ctx context.Context,
	restaurantID uuid.UUID,
	itemID uuid.UUID,
) (domain.MenuItem, error) {
	_, err := r.GetRestaurant(ctx, restaurantID)
	if err != nil {
		return domain.MenuItem{}, domain.ErrRestaurantNotFound
	}

	r.itemsMtx.RLock()
	defer r.itemsMtx.RUnlock()

	item, ok := r.items[restaurantID][itemID]
	if !ok {
		return domain.MenuItem{}, fmt.Errorf("not found: %w", domain.ErrMenuItemNotFound)
	}

	return itemDomainFromModel(itemID, item), nil
}

func (r *RestaurantsRepository) GetMenuItems(
	ctx context.Context,
	restaurantID uuid.UUID,
) ([]domain.MenuItem, error) {
	_, err := r.GetRestaurant(ctx, restaurantID)
	if err != nil {
		return []domain.MenuItem{}, domain.ErrRestaurantNotFound
	}

	r.itemsMtx.RLock()
	defer r.itemsMtx.RUnlock()

	keys := make([]uuid.UUID, 0, len(r.items[restaurantID]))
	for k := range r.items[restaurantID] {
		keys = append(keys, k)
	}

	// TODO: насколько норм такой способ упорядочить мапу, в которой ключ uuid.UUID?
	slices.SortFunc(keys, func(a, b uuid.UUID) int {
		return strings.Compare(a.String(), b.String())
	})

	menuItems := make([]domain.MenuItem, 0, len(keys))
	for _, v := range keys {
		menuItems = append(menuItems, itemDomainFromModel(v, r.items[restaurantID][v]))
	}

	return menuItems, nil
}

func (r *RestaurantsRepository) PutMenuItem(
	ctx context.Context,
	restaurantID uuid.UUID,
	itemDomain domain.MenuItem,
) (domain.MenuItem, error) {
	_, err := r.GetRestaurant(ctx, restaurantID)
	if err != nil {
		return domain.MenuItem{}, domain.ErrRestaurantNotFound
	}

	r.itemsMtx.Lock()
	defer r.itemsMtx.Unlock()

	itemID := itemDomain.ID
	r.items[restaurantID][itemID] = menuItemModelFromDomain(itemDomain)

	return itemDomain, nil
}

package memory

import (
	"context"
	"sync"
	"testing"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRestaurant() domain.Restaurant {
	return domain.NewRestaurant(uuid.New(), "KFC", true, 1000, "RUB")
}

func testMenuItem() domain.MenuItem {
	return domain.NewMenuItem(uuid.New(), "spicy chicken", "classic", 400, "RUB", true)
}

func TestRestaurantsRepository_CreateOrReplaceRestaurant_CreatesRestaurant(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	r, err := repo.GetRestaurant(ctx, restaurant.ID)
	require.NoError(t, err)

	assert.Equal(t, restaurant, r)
}

func TestRestaurantsRepository_CreateOrReplaceRestaurant_ReplaceWithoutDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	restaurants, err := repo.GetRestaurants(ctx)
	require.NoError(t, err)
	restaurantsCount := len(restaurants)

	changedRestaurant := domain.NewRestaurant(
		restaurant.ID,
		restaurant.Name,
		restaurant.AcceptingOrders,
		restaurant.MinimumOrderMinor,
		restaurant.Currency,
	)
	changedRestaurant.AcceptingOrders = false

	_, err = repo.CreateOrReplaceRestaurant(ctx, changedRestaurant)
	require.NoError(t, err)

	newRestaurants, err := repo.GetRestaurants(ctx)
	require.NoError(t, err)
	newRestaurantsCount := len(newRestaurants)
	assert.Equal(t, restaurantsCount, newRestaurantsCount)

	r, err := repo.GetRestaurant(ctx, restaurant.ID)
	require.NoError(t, err)

	assert.Equal(t, changedRestaurant, r)
}

func TestRestaurantsRepository_CreateOrReplaceRestaurant_ReplaceKeepsItems(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	_, err = repo.PutMenuItem(ctx, restaurant.ID, testMenuItem())
	require.NoError(t, err)
	_, err = repo.PutMenuItem(ctx, restaurant.ID, testMenuItem())
	require.NoError(t, err)

	items, err := repo.GetMenuItems(ctx, restaurant.ID)
	require.NoError(t, err)
	itemsCount := len(items)

	newr := domain.NewRestaurant(
		restaurant.ID,
		restaurant.Name,
		restaurant.AcceptingOrders,
		restaurant.MinimumOrderMinor,
		restaurant.Currency,
	)
	newr.AcceptingOrders = false
	repo.CreateOrReplaceRestaurant(ctx, newr)

	newItems, err := repo.GetMenuItems(ctx, restaurant.ID)
	require.NoError(t, err)
	newItemsCount := len(newItems)

	assert.Equal(t, itemsCount, newItemsCount)
	assert.Equal(t, items, newItems)
}

func TestRestaurantsRepository_PutMenuItem_CreatesItem(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	item := testMenuItem()
	_, err = repo.PutMenuItem(ctx, restaurant.ID, item)
	require.NoError(t, err)

	i, err := repo.GetMenuItem(ctx, restaurant.ID, item.ID)
	require.NoError(t, err)

	assert.Equal(t, item, i)
}

func TestRestaurantsRepository_GetMenuItem_IsolatedMenuItems(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant1 := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant1)
	require.NoError(t, err)

	restaurant2 := testRestaurant()
	_, err = repo.CreateOrReplaceRestaurant(ctx, restaurant2)
	require.NoError(t, err)

	i := testMenuItem()
	_, err = repo.PutMenuItem(ctx, restaurant1.ID, i)
	require.NoError(t, err)

	_, err = repo.GetMenuItem(ctx, restaurant2.ID, i.ID)
	require.ErrorIs(t, err, domain.ErrMenuItemNotFound)
}

func TestNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	tests := []struct {
		name    string
		f       func() error
		wantErr error
	}{
		{
			name: "unknown restaurant",
			f: func() error {
				_, err := repo.GetRestaurant(ctx, uuid.New())
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
		{
			name: "menu item of unknown restaurant",
			f: func() error {
				_, err := repo.GetMenuItem(ctx, uuid.New(), uuid.New())
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
		{
			name: "list items of unknown restaurant",
			f: func() error {
				_, err := repo.GetMenuItems(ctx, uuid.New())
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
		{
			name: "get unknown item of existing restaurant",
			f: func() error {
				_, err := repo.GetMenuItem(ctx, restaurant.ID, uuid.New())
				return err
			},
			wantErr: domain.ErrMenuItemNotFound,
		},
		{
			name: "put item to unknown restaurant",
			f: func() error {
				_, err := repo.PutMenuItem(ctx, uuid.New(), testMenuItem())
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.f(), tt.wantErr)
		})
	}
}

func TestRestaurantsRepository_GetMenuItems_Empty(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	items, err := repo.GetMenuItems(ctx, restaurant.ID)
	require.NoError(t, err)
	assert.Equal(t, []domain.MenuItem{}, items)
}

func TestRestaurantsRepository_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	workers := 50

	items := make([]domain.MenuItem, workers)
	for i, _ := range items {
		items[i] = testMenuItem()
	}

	wg := sync.WaitGroup{}

	for _, item := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
			assert.NoError(t, err)

			_, err = repo.PutMenuItem(ctx, restaurant.ID, item)
			assert.NoError(t, err)

			_, err = repo.GetMenuItems(ctx, restaurant.ID)
			assert.NoError(t, err)

			_, err = repo.GetRestaurants(ctx)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	res, err := repo.GetMenuItems(ctx, restaurant.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, items, res)
}

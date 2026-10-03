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

func TestCreateOrReplaceRestaurant_CreateNew(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	restaurant := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
	require.NoError(t, err)

	r, err := repo.GetRestaurant(ctx, restaurant.ID)
	require.NoError(t, err)

	assert.Equal(t, restaurant, r)
}

func TestCreateOrReplaceRestaurant_ReplaceKeepsData(t *testing.T) {
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

func TestCreateOrReplaceRestaurant_ReplaceKeepsMenu(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	r := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, r)
	require.NoError(t, err)

	_, err = repo.PutMenuItem(ctx, r.ID, testMenuItem())
	require.NoError(t, err)
	_, err = repo.PutMenuItem(ctx, r.ID, testMenuItem())
	require.NoError(t, err)

	items, err := repo.GetMenuItems(ctx, r.ID)
	require.NoError(t, err)
	itemsCount := len(items)

	newr := domain.NewRestaurant(
		r.ID,
		r.Name,
		r.AcceptingOrders,
		r.MinimumOrderMinor,
		r.Currency,
	)
	newr.AcceptingOrders = false
	repo.CreateOrReplaceRestaurant(ctx, newr)

	newItems, err := repo.GetMenuItems(ctx, r.ID)
	require.NoError(t, err)
	newItemsCount := len(newItems)

	assert.Equal(t, itemsCount, newItemsCount)
	assert.Equal(t, items, newItems)
}

func TestPutMenuItem_CreateNew(t *testing.T) {
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

func TestIsolatedMenuItems(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	r1 := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, r1)
	require.NoError(t, err)

	r2 := testRestaurant()
	_, err = repo.CreateOrReplaceRestaurant(ctx, r2)
	require.NoError(t, err)

	i := testMenuItem()
	_, err = repo.PutMenuItem(ctx, r1.ID, i)
	require.NoError(t, err)

	_, err = repo.GetMenuItem(ctx, r2.ID, i.ID)
	require.ErrorIs(t, err, domain.ErrMenuItemNotFound)
}

func TestNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	r := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, r)
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
				_, err := repo.GetMenuItem(ctx, r.ID, uuid.New())
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

func TestGetMenuItems_NoItems(t *testing.T) {
	ctx := context.Background()
	repo := NewRestaurantsRepository()

	r := testRestaurant()
	_, err := repo.CreateOrReplaceRestaurant(ctx, r)
	require.NoError(t, err)

	items, err := repo.GetMenuItems(ctx, r.ID)
	require.NoError(t, err)
	assert.Equal(t, []domain.MenuItem{}, items)
}

func TestRepo_ConcurrentAcces(t *testing.T) {
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

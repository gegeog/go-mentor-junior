package service

import (
	"context"
	"testing"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestaurantsService_PutRestaurant_OK(t *testing.T) {
	repo := NewRestaurantsRepositoryMock(t)

	restaurant := domain.NewRestaurant(uuid.New(), "cafe", true, 1000, "USD")
	repo.GetRestaurantMock.
		Expect(minimock.AnyContext, restaurant.ID).
		Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)

	repo.CreateOrReplaceRestaurantMock.
		Expect(minimock.AnyContext, restaurant).
		Return(restaurant, nil)

	svc := NewRestaurantsService(repo)
	got, err := svc.PutRestaurant(context.Background(), restaurant)

	require.NoError(t, err)
	assert.Equal(t, restaurant, got)
}

func TestRestaurantsService_PutRestaurant_CurrencyMismatch(t *testing.T) {
	repo := NewRestaurantsRepositoryMock(t)

	restaurantID := uuid.New()
	existed := domain.NewRestaurant(restaurantID, "cafe", true, 1000, "USD")
	updated := domain.NewRestaurant(restaurantID, "cafe", true, 1000, "RUB")

	repo.GetRestaurantMock.
		Expect(minimock.AnyContext, existed.ID).
		Return(existed, nil)

	svc := NewRestaurantsService(repo)
	_, err := svc.PutRestaurant(context.Background(), updated)
	require.ErrorIs(t, err, domain.ErrCurrencyMismatch)
}

func TestRestaurantsService_PutMenuItem_OK(t *testing.T) {
	repo := NewRestaurantsRepositoryMock(t)

	restaurantID := uuid.New()
	restaurant := domain.NewRestaurant(restaurantID, "kfc", true, 1000, "USD")
	item := domain.NewMenuItem(uuid.New(), "wings", "spicy", 1000, "USD", true)

	repo.GetRestaurantMock.
		Expect(minimock.AnyContext, restaurantID).
		Return(restaurant, nil)

	repo.PutMenuItemMock.
		Expect(minimock.AnyContext, restaurantID, item).
		Return(item, nil)

	svc := NewRestaurantsService(repo)
	got, err := svc.PutMenuItem(context.Background(), restaurantID, item)
	require.NoError(t, err)
	assert.Equal(t, item, got)
}

func TestRestaurantsService_PutMenuItem_Errors(t *testing.T) {
	restaurantID := uuid.New()
	item := domain.NewMenuItem(uuid.New(), "wings", "spicy", 200, "RUB", true)

	tests := []struct {
		name    string
		f       func(*RestaurantsRepositoryMock)
		wantErr error
	}{
		{
			name: "restaurant not found",
			f: func(m *RestaurantsRepositoryMock) {
				m.GetRestaurantMock.Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
		{
			name: "currency mismatch",
			f: func(m *RestaurantsRepositoryMock) {
				m.GetRestaurantMock.Return(domain.NewRestaurant(restaurantID, "rest", true, 1000, "USD"), nil)
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRestaurantsRepositoryMock(t)
			tt.f(repo)

			svc := NewRestaurantsService(repo)
			_, err := svc.PutMenuItem(context.Background(), restaurantID, item)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRestaurantsService_Getters_OK(t *testing.T) {
	ctx := context.Background()
	restaurantID := uuid.New()
	restaurant := domain.NewRestaurant(restaurantID, "mcdonalds", true, 1000, "USD")
	item := domain.NewMenuItem(uuid.New(), "wings", "spicy", 500, "USD", true)
	restaurants := []domain.Restaurant{
		domain.NewRestaurant(uuid.New(), "a", true, 100, "USD"),
		domain.NewRestaurant(uuid.New(), "b", true, 200, "EUR"),
	}
	items := []domain.MenuItem{item}

	tests := []struct {
		name   string
		setup  func(m *RestaurantsRepositoryMock)
		assert func(t *testing.T, s *RestaurantsService)
	}{
		{
			name: "get restaurant",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetRestaurantMock.
					Expect(minimock.AnyContext, restaurantID).
					Return(restaurant, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetRestaurant(ctx, restaurantID)
				require.NoError(t, err)
				assert.Equal(t, restaurant, got)
			},
		},
		{
			name: "get restaurants",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetRestaurantsMock.Return(restaurants, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetRestaurants(ctx)
				require.NoError(t, err)
				assert.Equal(t, restaurants, got)
			},
		},
		{
			name: "get restaurants empty",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetRestaurantsMock.Return([]domain.Restaurant{}, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetRestaurants(ctx)
				require.NoError(t, err)
				assert.Empty(t, got)
			},
		},
		{
			name: "get item",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetMenuItemMock.
					Expect(minimock.AnyContext, restaurantID, item.ID).
					Return(item, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetItem(ctx, restaurantID, item.ID)
				require.NoError(t, err)
				assert.Equal(t, item, got)
			},
		},
		{
			name: "get items",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetMenuItemsMock.
					Expect(minimock.AnyContext, restaurantID).
					Return(items, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetItems(ctx, restaurantID)
				require.NoError(t, err)
				assert.Equal(t, items, got)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRestaurantsRepositoryMock(t)
			tt.setup(repo)
			tt.assert(t, NewRestaurantsService(repo))
		})
	}
}

func TestRestaurantsService_Getters_RepoErrors(t *testing.T) {
	ctx := context.Background()
	restaurantID := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name    string
		setup   func(m *RestaurantsRepositoryMock)
		call    func(s *RestaurantsService) error
		wantErr error
	}{
		{
			name: "get restaurant storage error",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetRestaurantMock.Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)
			},
			call: func(s *RestaurantsService) error {
				_, err := s.GetRestaurant(ctx, restaurantID)
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
		{
			name: "get menu item not found storage error",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetMenuItemMock.Return(domain.MenuItem{}, domain.ErrMenuItemNotFound)
			},
			call: func(s *RestaurantsService) error {
				_, err := s.GetItem(ctx, restaurantID, itemID)
				return err
			},
			wantErr: domain.ErrMenuItemNotFound,
		},
		{
			name: "get menu item restaurant not found storage error",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetMenuItemMock.Return(domain.MenuItem{}, domain.ErrRestaurantNotFound)
			},
			call: func(s *RestaurantsService) error {
				_, err := s.GetItem(ctx, restaurantID, itemID)
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
		{
			name: "get menu items restaurant not found storage error",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetMenuItemsMock.Return(nil, domain.ErrRestaurantNotFound)
			},
			call: func(s *RestaurantsService) error {
				_, err := s.GetItems(ctx, restaurantID)
				return err
			},
			wantErr: domain.ErrRestaurantNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRestaurantsRepositoryMock(t)
			tt.setup(repo)

			err := tt.call(NewRestaurantsService(repo))
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

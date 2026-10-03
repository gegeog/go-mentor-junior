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

func TestPutRestaurant_OK(t *testing.T) {
	repo := NewRestaurantsRepositoryMock(t)

	r := domain.NewRestaurant(uuid.New(), "cafe", true, 1000, "USD")
	repo.GetRestaurantMock.
		Expect(minimock.AnyContext, r.ID).
		Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)

	repo.CreateOrReplaceRestaurantMock.
		Expect(minimock.AnyContext, r).
		Return(r, nil)

	svc := NewRestaurantsService(repo)
	got, err := svc.PutRestaurant(context.Background(), r)

	require.NoError(t, err)
	assert.Equal(t, r, got)
}

func TestPutRestaurant_CurrencyMismatch(t *testing.T) {
	repo := NewRestaurantsRepositoryMock(t)

	rid := uuid.New()
	existed := domain.NewRestaurant(rid, "cafe", true, 1000, "USD")
	updated := domain.NewRestaurant(rid, "cafe", true, 1000, "RUB")

	repo.GetRestaurantMock.
		Expect(minimock.AnyContext, existed.ID).
		Return(existed, nil)

	svc := NewRestaurantsService(repo)
	_, err := svc.PutRestaurant(context.Background(), updated)
	require.ErrorIs(t, err, domain.ErrCurrencyMismatch)
}

func TestPutMenuItem_OK(t *testing.T) {
	repo := NewRestaurantsRepositoryMock(t)

	rid := uuid.New()
	r := domain.NewRestaurant(rid, "kfc", true, 1000, "USD")
	item := domain.NewMenuItem(uuid.New(), "wings", "spicy", 1000, "USD", true)

	repo.GetRestaurantMock.
		Expect(minimock.AnyContext, rid).
		Return(r, nil)

	repo.PutMenuItemMock.
		Expect(minimock.AnyContext, rid, item).
		Return(item, nil)

	svc := NewRestaurantsService(repo)
	got, err := svc.PutMenuItem(context.Background(), rid, item)
	require.NoError(t, err)
	assert.Equal(t, item, got)
}

func TestPutMenuItem_Errors(t *testing.T) {
	rid := uuid.New()
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
				m.GetRestaurantMock.Return(domain.NewRestaurant(rid, "rest", true, 1000, "USD"), nil)
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRestaurantsRepositoryMock(t)
			tt.f(repo)

			svc := NewRestaurantsService(repo)
			_, err := svc.PutMenuItem(context.Background(), rid, item)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestServiceGetters_OK(t *testing.T) {
	ctx := context.Background()
	rid := uuid.New()
	restaurant := domain.NewRestaurant(rid, "mcdonalds", true, 1000, "USD")
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
					Expect(minimock.AnyContext, rid).
					Return(restaurant, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetRestaurant(ctx, rid)
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
					Expect(minimock.AnyContext, rid, item.ID).
					Return(item, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetItem(ctx, rid, item.ID)
				require.NoError(t, err)
				assert.Equal(t, item, got)
			},
		},
		{
			name: "get items",
			setup: func(m *RestaurantsRepositoryMock) {
				m.GetMenuItemsMock.
					Expect(minimock.AnyContext, rid).
					Return(items, nil)
			},
			assert: func(t *testing.T, s *RestaurantsService) {
				got, err := s.GetItems(ctx, rid)
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

func TestServiceGetters_RepoErrors(t *testing.T) {
	ctx := context.Background()
	rid := uuid.New()
	iid := uuid.New()

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
				_, err := s.GetRestaurant(ctx, rid)
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
				_, err := s.GetItem(ctx, rid, iid)
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
				_, err := s.GetItem(ctx, rid, iid)
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
				_, err := s.GetItems(ctx, rid)
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

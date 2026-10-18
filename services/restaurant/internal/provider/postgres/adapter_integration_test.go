package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

func TestRepo(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18.6-bookworm",
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			fmt.Println("failed to terminate container")
		}
	}()

	connStr, err := pgContainer.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	require.NoError(t, err)

	conncfg, err := pgxpool.ParseConfig(connStr)
	require.NoError(t, err)

	pgxpool, err := pgxpool.NewWithConfig(ctx, conncfg)
	require.NoError(t, err)
	defer pgxpool.Close()

	err = migrations.Run(*conncfg.ConnConfig, zap.NewNop(), "migrate", "up")
	if err != nil {
		t.Fatal(err)
	}

	restaurantID := uuid.New()

	repo := NewRestaurantsRepository(pgxpool)
	t.Run("create restaurant", func(t *testing.T) {
		restaurant := domain.NewRestaurant(restaurantID, "cafeteria", true, 5000, "RUB")
		_, err := repo.CreateOrReplaceRestaurant(ctx, restaurant)
		require.NoError(t, err)

		got, err := repo.GetRestaurant(ctx, restaurant.ID)
		require.NoError(t, err)

		require.Equal(t, restaurant, got)
	})

	t.Run("replace restaurant", func(t *testing.T) {
		r, err := repo.GetRestaurant(ctx, restaurantID)
		require.NoError(t, err)

		_, err = repo.CreateOrReplaceRestaurant(
			ctx, domain.NewRestaurant(restaurantID, "cafeteria", false, 5000, "RUB"),
		)
		require.NoError(t, err)

		updr, err := repo.GetRestaurant(ctx, restaurantID)
		require.NoError(t, err)

		require.Equal(t, r.Name, updr.Name)
		require.Equal(t, r.MinimumOrderMinor, updr.MinimumOrderMinor)
		require.Equal(t, r.Currency, updr.Currency)

		require.NotEqual(t, r.AcceptingOrders, updr.AcceptingOrders)
	})

	t.Run("change restaurant currency", func(t *testing.T) {
		r, err := repo.GetRestaurant(ctx, restaurantID)
		require.NoError(t, err)

		_, err = repo.CreateOrReplaceRestaurant(
			ctx, domain.NewRestaurant(restaurantID, "cafeteria", false, 5000, "USD"),
		)
		if !errors.Is(err, domain.ErrCurrencyMismatch) {
			t.Fatalf("changing currency of existed restaurant")
		}

		updr, err := repo.GetRestaurant(ctx, restaurantID)
		require.NoError(t, err)

		require.Equal(t, r.Name, updr.Name)
		require.Equal(t, r.MinimumOrderMinor, updr.MinimumOrderMinor)
		require.Equal(t, r.Currency, updr.Currency)
		require.Equal(t, r.AcceptingOrders, updr.AcceptingOrders)
	})

	itemID := uuid.New()
	t.Run("create menu item", func(t *testing.T) {
		item := domain.NewMenuItem(itemID, "soup", "", 4000, "RUB", true)
		_, err := repo.PutMenuItem(ctx, restaurantID, item)
		require.NoError(t, err)

		got, err := repo.GetMenuItem(ctx, restaurantID, item.ID)
		require.NoError(t, err)
		require.Equal(t, item, got)
	})

	t.Run("replace menu item", func(t *testing.T) {
		i, err := repo.GetMenuItem(ctx, restaurantID, itemID)
		require.NoError(t, err)

		item := domain.NewMenuItem(itemID, "soup", "", 4000, "RUB", false)
		_, err = repo.PutMenuItem(ctx, restaurantID, item)
		require.NoError(t, err)

		updi, err := repo.GetMenuItem(ctx, restaurantID, itemID)
		require.NoError(t, err)

		require.Equal(t, i.Name, updi.Name)
		require.Equal(t, i.Description, updi.Description)
		require.Equal(t, i.PriceMinor, updi.PriceMinor)
		require.Equal(t, i.Currency, updi.Currency)
		require.NotEqual(t, i.Available, updi.Available)
	})

	t.Run("change item currency", func(t *testing.T) {
		i, err := repo.GetMenuItem(ctx, restaurantID, itemID)
		require.NoError(t, err)

		_, err = repo.PutMenuItem(
			ctx, restaurantID, domain.NewMenuItem(itemID, "pizza", "cheesy", 14000, "USD", false),
		)
		if !errors.Is(err, domain.ErrCurrencyMismatch) {
			t.Fatalf("changing currency of existed item")
		}

		updi, err := repo.GetMenuItem(ctx, restaurantID, itemID)
		require.NoError(t, err)

		require.Equal(t, i.Name, updi.Name)
		require.Equal(t, i.Description, updi.Description)
		require.Equal(t, i.PriceMinor, updi.PriceMinor)
		require.Equal(t, i.Currency, updi.Currency)
		require.Equal(t, i.Available, updi.Available)
	})

	t.Run("menu item belongs to one restaurant", func(t *testing.T) {
		owner := domain.NewRestaurant(uuid.New(), "owner", true, 500, "RUB")
		another := domain.NewRestaurant(uuid.New(), "another", true, 500, "RUB")
		item := domain.NewMenuItem(uuid.New(), "soup", "hot", 400, "RUB", true)

		_, err := repo.CreateOrReplaceRestaurant(ctx, owner)
		require.NoError(t, err)

		_, err = repo.CreateOrReplaceRestaurant(ctx, another)
		require.NoError(t, err)

		_, err = repo.PutMenuItem(ctx, owner.ID, item)
		require.NoError(t, err)

		_, err = repo.GetMenuItem(ctx, another.ID, item.ID)
		require.ErrorIs(t, err, domain.ErrMenuItemNotFound)

		_, err = repo.PutMenuItem(ctx, another.ID, item)
		require.ErrorIs(t, err, domain.ErrInvalidRestaurant)

		got, err := repo.GetMenuItem(ctx, owner.ID, item.ID)
		require.NoError(t, err)
		require.Equal(t, item, got)

		anotherItems, err := repo.GetMenuItems(ctx, another.ID)
		require.NoError(t, err)
		require.Empty(t, anotherItems)
	})

	t.Run("no records", func(t *testing.T) {
		rid := uuid.New()
		iid := uuid.New()

		_, err := repo.GetRestaurant(ctx, rid)
		require.ErrorIs(t, err, domain.ErrRestaurantNotFound)

		_, err = repo.GetMenuItems(ctx, rid)
		require.ErrorIs(t, err, domain.ErrRestaurantNotFound)

		_, err = repo.GetMenuItem(ctx, rid, iid)
		require.ErrorIs(t, err, domain.ErrRestaurantNotFound)

		_, err = repo.CreateOrReplaceRestaurant(
			ctx,
			domain.NewRestaurant(rid, "t", true, 500, "rub"),
		)
		require.NoError(t, err)

		_, err = repo.GetMenuItem(ctx, rid, iid)
		require.ErrorIs(t, err, domain.ErrMenuItemNotFound)

		items, err := repo.GetMenuItems(ctx, rid)
		require.NoError(t, err)
		require.Empty(t, items)
	})

	t.Run("wrong first item currency", func(t *testing.T) {
		r := domain.NewRestaurant(uuid.New(), "pizza ru", true, 4000, "RUB")
		i := domain.NewMenuItem(uuid.New(), "chees", "", 200, "USD", false)

		_, err := repo.CreateOrReplaceRestaurant(ctx, r)
		require.NoError(t, err)

		_, err = repo.PutMenuItem(ctx, r.ID, i)
		require.ErrorIs(t, err, domain.ErrCurrencyMismatch)

		_, err = repo.GetMenuItem(ctx, r.ID, i.ID)
		require.ErrorIs(t, err, domain.ErrMenuItemNotFound)
	})

	t.Run("concurrent restaurants with conflict id", func(t *testing.T) {
		id := uuid.New()
		restaurants := []domain.Restaurant{
			{
				ID:                id,
				Name:              "pizza.ru",
				AcceptingOrders:   true,
				MinimumOrderMinor: 3000,
				Currency:          "RUB",
			},
			{
				ID:                id,
				Name:              "burgerlab",
				AcceptingOrders:   true,
				MinimumOrderMinor: 4000,
				Currency:          "USD",
			},
		}

		start := make(chan struct{})
		wg := sync.WaitGroup{}

		var errs [2]error
		for i, r := range restaurants {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				_, err := repo.CreateOrReplaceRestaurant(ctx, r)
				errs[i] = err
			}()
		}

		close(start)
		wg.Wait()

		got, err := repo.GetRestaurant(ctx, id)
		require.NoError(t, err)

		if errs[0] == nil {
			require.ErrorIs(t, errs[1], domain.ErrCurrencyMismatch)
			require.Equal(t, restaurants[0], got)
		} else {
			require.ErrorIs(t, errs[0], domain.ErrCurrencyMismatch)
			require.NoError(t, errs[1])
			require.Equal(t, restaurants[1], got)
		}

		require.Contains(t, restaurants, got)
	})

	t.Run("concurrent adding menu item", func(t *testing.T) {
		restaurants := []domain.Restaurant{
			{
				ID:                uuid.New(),
				Name:              "pizza.ru",
				AcceptingOrders:   true,
				MinimumOrderMinor: 3000,
				Currency:          "RUB",
			},
			{
				ID:                uuid.New(),
				Name:              "burgerlab",
				AcceptingOrders:   true,
				MinimumOrderMinor: 4000,
				Currency:          "RUB",
			},
		}

		for _, v := range restaurants {
			_, err := repo.CreateOrReplaceRestaurant(ctx, v)
			require.NoError(t, err)
		}

		item := domain.NewMenuItem(uuid.New(), "soup", "tomato", 5000, "RUB", true)

		start := make(chan struct{})
		wg := sync.WaitGroup{}
		var errs [2]error

		for i, r := range restaurants {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				_, err := repo.PutMenuItem(ctx, r.ID, item)
				errs[i] = err
			}()
		}

		close(start)
		wg.Wait()

		winner := 0
		if errs[0] == nil {
			require.ErrorIs(t, errs[1], domain.ErrInvalidRestaurant)
		} else {
			require.ErrorIs(t, errs[0], domain.ErrInvalidRestaurant)
			require.NoError(t, errs[1])
			winner = 1
		}

		got, err := repo.GetMenuItem(ctx, restaurants[winner].ID, item.ID)
		require.NoError(t, err)
		require.Equal(t, item, got)

		loser := 1 - winner
		_, err = repo.GetMenuItem(ctx, restaurants[loser].ID, item.ID)
		require.ErrorIs(t, err, domain.ErrMenuItemNotFound)
	})
}

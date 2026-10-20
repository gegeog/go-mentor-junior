package postgres

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantsRepository struct {
	queries *Queries
}

func NewRestaurantsRepository(pool *pgxpool.Pool) *RestaurantsRepository {
	return &RestaurantsRepository{
		queries: New(pool),
	}
}

func (r *RestaurantsRepository) GetMenuItem(
	ctx context.Context,
	restaurantID uuid.UUID,
	itemID uuid.UUID,
) (domain.MenuItem, error) {
	ps := GetMenuItemParams{
		RestaurantID: pgtype.UUID{
			Bytes: restaurantID,
			Valid: true,
		},
		ID: pgtype.UUID{
			Bytes: itemID,
			Valid: true,
		},
	}

	if _, err := r.queries.GetRestaurant(ctx, ps.RestaurantID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MenuItem{}, fmt.Errorf("restaurant not found: %w", domain.ErrRestaurantNotFound)
		}
		return domain.MenuItem{}, fmt.Errorf("get restaurant: %s: %w", ps.RestaurantID, err)
	}

	row, err := r.queries.GetMenuItem(ctx, ps)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MenuItem{}, fmt.Errorf("item not found: %w", domain.ErrMenuItemNotFound)
		}
		return domain.MenuItem{}, fmt.Errorf("get menu item: %s: %w", itemID, err)
	}

	return domain.NewMenuItem(
		row.ID.Bytes,
		row.Name,
		row.Description,
		row.PriceMinor,
		row.Currency,
		row.Available,
	), nil
}

func (r *RestaurantsRepository) GetMenuItems(
	ctx context.Context,
	restaurantID uuid.UUID,
) ([]domain.MenuItem, error) {
	uuid := pgtype.UUID{
		Bytes: restaurantID,
		Valid: true,
	}

	if _, err := r.queries.GetRestaurant(ctx, uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.MenuItem{}, fmt.Errorf("get restaurant: %s: %w", restaurantID, domain.ErrRestaurantNotFound)
		}
		return []domain.MenuItem{}, fmt.Errorf("get restaurant: %s: %w", restaurantID, err)
	}

	rows, err := r.queries.GetMenuItems(
		ctx,
		uuid,
	)

	if err != nil {
		return []domain.MenuItem{}, fmt.Errorf(
			"get menu items of restaurant: %s: %w",
			restaurantID,
			err,
		)
	}

	items := make([]domain.MenuItem, len(rows))
	for i, v := range rows {
		items[i] = domain.NewMenuItem(
			v.ID.Bytes,
			v.Name,
			v.Description,
			v.PriceMinor,
			v.Currency,
			v.Available,
		)
	}

	return items, nil
}

func (r *RestaurantsRepository) PutMenuItem(
	ctx context.Context,
	restaurantID uuid.UUID,
	itemDomain domain.MenuItem,
) (domain.MenuItem, error) {
	item, err := r.queries.PutMenuItem(ctx, PutMenuItemParams{
		ID: pgtype.UUID{
			Bytes: itemDomain.ID,
			Valid: true,
		},
		Name:        itemDomain.Name,
		Description: itemDomain.Description,
		PriceMinor:  itemDomain.PriceMinor,
		Currency:    itemDomain.Currency,
		Available:   itemDomain.Available,
		ID_2: pgtype.UUID{
			Bytes: restaurantID,
			Valid: true,
		},
	})

	// todo: нужно ли такую большую проверку делать на уровне репо?
	// если я уже по сути делаю её на уровне сервиса
	// я как будто оверинжинирю тут..
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			restaurant, restaurantErr := r.GetRestaurant(ctx, restaurantID)

			if restaurantErr != nil {
				return domain.MenuItem{}, fmt.Errorf(
					"put menu item %s: %w", itemDomain.ID, restaurantErr,
				)
			}
			if restaurant.Currency != itemDomain.Currency {
				return domain.MenuItem{}, fmt.Errorf(
					"put menu item %s: %w", itemDomain.ID, domain.ErrCurrencyMismatch,
				)
			}

			_, getErr := r.queries.GetMenuItem(ctx, GetMenuItemParams{
				RestaurantID: pgtype.UUID{Bytes: restaurantID, Valid: true},
				ID:           pgtype.UUID{Bytes: itemDomain.ID, Valid: true},
			})

			if errors.Is(getErr, pgx.ErrNoRows) {
				return domain.MenuItem{}, fmt.Errorf(
					"put menu item %s: %w",
					itemDomain.ID,
					domain.ErrInvalidRestaurant,
				)
			}
			if getErr != nil {
				return domain.MenuItem{}, fmt.Errorf(
					"get menu item %s after failed put: %w",
					itemDomain.ID,
					getErr,
				)
			}

			return domain.MenuItem{}, fmt.Errorf(
				"put menu item %s: %w",
				itemDomain.ID,
				domain.ErrCurrencyMismatch,
			)
		}

		return domain.MenuItem{}, fmt.Errorf("put item: %s: %w", itemDomain.ID, err)
	}

	return domain.NewMenuItem(
		item.ID.Bytes,
		item.Name,
		item.Description,
		item.PriceMinor,
		item.Currency,
		item.Available,
	), nil
}

func (r *RestaurantsRepository) GetRestaurant(
	ctx context.Context,
	id uuid.UUID,
) (domain.Restaurant, error) {
	restaurant, err := r.queries.GetRestaurant(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Restaurant{}, fmt.Errorf("get restaurant: %s: %w", id, domain.ErrRestaurantNotFound)
		}
		return domain.Restaurant{}, fmt.Errorf("get restaurant: %s: %w", id, err)
	}

	return domain.NewRestaurant(
		restaurant.ID.Bytes,
		restaurant.Name,
		restaurant.AcceptingOrders,
		restaurant.MinimumOrderMinor,
		restaurant.Currency,
	), nil
}

func (r *RestaurantsRepository) GetRestaurants(
	ctx context.Context,
) ([]domain.Restaurant, error) {
	restaurants, err := r.queries.GetRestaurants(ctx)
	if err != nil {
		return []domain.Restaurant{}, fmt.Errorf("get restaurants: %w", err)
	}

	restaurantsDomains := make([]domain.Restaurant, len(restaurants))
	for k, v := range restaurants {
		restaurantsDomains[k] = domain.NewRestaurant(
			v.ID.Bytes,
			v.Name,
			v.AcceptingOrders,
			v.MinimumOrderMinor,
			v.Currency,
		)
	}

	return restaurantsDomains, nil
}

func (r *RestaurantsRepository) CreateOrReplaceRestaurant(
	ctx context.Context,
	restaurant domain.Restaurant,
) (domain.Restaurant, error) {
	restaurantNew, err := r.queries.CreateOrReplaceRestaurant(
		ctx,
		CreateOrReplaceRestaurantParams{
			ID: pgtype.UUID{
				Bytes: restaurant.ID,
				Valid: true,
			},
			Name:              restaurant.Name,
			AcceptingOrders:   restaurant.AcceptingOrders,
			MinimumOrderMinor: restaurant.MinimumOrderMinor,
			Currency:          restaurant.Currency,
		},
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Restaurant{}, fmt.Errorf("put restaurant: %s: %w", restaurant.ID, domain.ErrCurrencyMismatch)
		}
		return domain.Restaurant{}, fmt.Errorf("put restaurant: %s: %w", restaurant.ID, err)
	}

	return domain.NewRestaurant(
		restaurantNew.ID.Bytes,
		restaurantNew.Name,
		restaurantNew.AcceptingOrders,
		restaurantNew.MinimumOrderMinor,
		restaurantNew.Currency,
	), nil
}

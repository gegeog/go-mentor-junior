-- name: GetRestaurants :many
SELECT id, name, accepting_orders, minimum_order_minor, currency
FROM restaurants.restaurant
ORDER BY id;

-- name: GetRestaurant :one
SELECT id, name, accepting_orders, minimum_order_minor, currency
FROM restaurants.restaurant
WHERE id=$1;

-- name: CreateOrReplaceRestaurant :one
INSERT INTO restaurants.restaurant (id, name, accepting_orders, minimum_order_minor, currency)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  accepting_orders = EXCLUDED.accepting_orders,
  minimum_order_minor = EXCLUDED.minimum_order_minor
WHERE restaurants.restaurant.currency = EXCLUDED.currency
RETURNING id, name, accepting_orders, minimum_order_minor, currency;

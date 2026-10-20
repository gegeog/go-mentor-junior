-- name: GetMenuItem :one
SELECT id, name, description, price_minor, currency, available
FROM restaurants.menu_item
WHERE restaurant_id = $1 AND id = $2;

-- name: GetMenuItems :many
SELECT id, name, description, price_minor, currency, available
FROM restaurants.menu_item
WHERE restaurant_id=$1
ORDER BY id;

-- name: PutMenuItem :one
INSERT INTO restaurants.menu_item (
    id, name, description, price_minor,
    currency, available, restaurant_id
)
SELECT $1, $2, $3, $4, $5, $6, restaurant.id
FROM restaurants.restaurant AS restaurant
WHERE restaurant.id = $7
  AND restaurant.currency = $5
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price_minor = EXCLUDED.price_minor,
    available = EXCLUDED.available
WHERE restaurants.menu_item.restaurant_id = EXCLUDED.restaurant_id
AND restaurants.menu_item.currency = EXCLUDED.currency
RETURNING id, name, description, price_minor, currency, available, restaurant_id;

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
INSERT INTO restaurants.menu_item
  (id, name, description, price_minor, currency, available, restaurant_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  price_minor = EXCLUDED.price_minor,
  available = EXCLUDED.available
WHERE restaurants.menu_item.restaurant_id = EXCLUDED.restaurant_id
RETURNING
  id, name, description, price_minor, currency, available, restaurant_id;

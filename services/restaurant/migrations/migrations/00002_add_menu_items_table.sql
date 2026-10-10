-- +goose Up
CREATE TABLE IF NOT EXISTS restaurants.menu_item (
  id                  UUID            NOT NULL PRIMARY KEY,
  name                TEXT            NOT NULL CHECK (char_length(name) > 0),
  description         TEXT            NOT NULL,
  price_minor         BIGINT          NOT NULL CHECK (price_minor > 0),
  currency            TEXT            NOT NULL CHECK (char_length(currency) > 0),
  available           BOOLEAN         NOT NULL,
  restaurant_id       UUID            NOT NULL REFERENCES restaurants.restaurant (id)
);

CREATE INDEX IF NOT EXISTS idx_menu_items_restaurant_id
    ON restaurants.menu_item USING BTREE (restaurant_id);

-- +goose Down
DROP INDEX IF EXISTS restaurants.idx_menu_items_restaurant_id;
DROP TABLE IF EXISTS restaurants.menu_item;

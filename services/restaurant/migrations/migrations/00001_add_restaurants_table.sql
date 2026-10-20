-- +goose Up
CREATE SCHEMA IF NOT EXISTS restaurants;

CREATE TABLE IF NOT EXISTS restaurants.restaurant (
  id                  UUID            NOT NULL PRIMARY KEY,
  name                TEXT            NOT NULL CHECK (char_length(name) > 0),
  accepting_orders    BOOLEAN         NOT NULL,
  minimum_order_minor BIGINT          NOT NULL CHECK (minimum_order_minor > 0),
  currency            TEXT            NOT NULL CHECK (char_length(currency) > 0)
);

-- +goose Down
DROP TABLE IF EXISTS restaurants.restaurant;
DROP SCHEMA IF EXISTS restaurants;

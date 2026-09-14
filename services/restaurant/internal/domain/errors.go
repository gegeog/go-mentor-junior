package domain

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")

	ErrRestaurantNotFound     = errors.New("restaurant not found")
	ErrRestaurantAccessDenied = errors.New("restaurant access denied")
	ErrInvalidRestaurant      = errors.New("invalid restaurant")

	ErrMenuItemNotFound = errors.New("menu item not found")
	ErrInvalidMenuItem  = errors.New("invalid menu item")

	ErrCurrencyMismatch = errors.New("currency mismatch")
)

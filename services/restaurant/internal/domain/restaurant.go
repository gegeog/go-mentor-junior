package domain

import (
	"fmt"
)

type Restaurant struct {
	ID                 string
	Name               string
	AcceptingOrders    bool
	MinimumOrderMminor int64
	Currency           string
}

func NewRestaurant(
	id, name string,
	acceptingOrders bool,
	minimumOrderMinor int64,
	currency string,
) Restaurant {
	return Restaurant{
		ID:                 id,
		Name:               name,
		AcceptingOrders:    acceptingOrders,
		MinimumOrderMminor: minimumOrderMinor,
		Currency:           currency,
	}
}

func (r Restaurant) Validate() error {
	if r.Name == "" {
		return fmt.Errorf(
			"invalid restaurant name: %w",
			ErrInvalidRestaurant,
		)
	}

	if r.MinimumOrderMminor <= 0 {
		return fmt.Errorf(
			"minimum order amount less or equal to zero: %w",
			ErrInvalidRestaurant,
		)
	}

	return nil
}

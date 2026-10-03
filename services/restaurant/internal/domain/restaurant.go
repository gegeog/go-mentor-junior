package domain

import (
	"uuid"
)

type Restaurant struct {
	ID                uuid.UUID
	Name              string
	AcceptingOrders   bool
	MinimumOrderMinor int64
	Currency          string
}

func NewRestaurant(
	id uuid.UUID,
	name string,
	acceptingOrders bool,
	minimumOrderMinor int64,
	currency string,
) Restaurant {
	return Restaurant{
		ID:                id,
		Name:              name,
		AcceptingOrders:   acceptingOrders,
		MinimumOrderMinor: minimumOrderMinor,
		Currency:          currency,
	}
}

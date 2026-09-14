package domain

import (
	"fmt"
)

type MenuItem struct {
	ID          string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	Available   bool
}

func NewMenuItem(
	id, name, description string,
	price int64,
	currency string,
	available bool,
) MenuItem {
	return MenuItem{
		ID:          id,
		Name:        name,
		Description: description,
		PriceMinor:  price,
		Currency:    currency,
		Available:   available,
	}
}

func (i MenuItem) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("invalid item id: %w", ErrInvalidMenuItem)
	}

	if i.Name == "" {
		return fmt.Errorf("invalid item name: %w", ErrInvalidMenuItem)
	}

	if i.PriceMinor < 0 {
		return fmt.Errorf("invalid item price: %w", ErrInvalidMenuItem)
	}

	return nil
}

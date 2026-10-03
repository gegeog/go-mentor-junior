package domain

import (
	"uuid"
)

type MenuItem struct {
	ID          uuid.UUID
	Name        string
	Description string
	//TODO: касательно точности...вроде нам будет достаточно инта, типа мы не переполним его позициями из меню??
	// ну и вычисления над целыми числами проводим, а отображем десятичные.
	// в общем я не вижу проблем при использовании plain int64, подскажи плиз
	PriceMinor int64
	Currency   string
	Available  bool
}

func NewMenuItem(
	id uuid.UUID,
	name string,
	description string,
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

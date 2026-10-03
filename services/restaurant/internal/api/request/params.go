package request

import (
	"fmt"
	"net/http"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
)

const (
	PathParamRestaurantID = "restaurant_id"
	PathParamMenuItemID   = "menu_item_id"

	HeaderRestaurantID = "X-Restaurant-ID"
)

func ParseRestaurantID(r *http.Request) (uuid.UUID, error) {
	return parsePathUUID(r, PathParamRestaurantID, domain.ErrInvalidRestaurant)
}

func ParseMenuItemID(r *http.Request) (uuid.UUID, error) {
	return parsePathUUID(r, PathParamMenuItemID, domain.ErrInvalidMenuItem)
}

func parsePathUUID(r *http.Request, param string, domainErr error) (uuid.UUID, error) {
	raw := r.PathValue(param)
	if raw == "" {
		return uuid.Nil(), fmt.Errorf("missing path param %s: %w", param, domainErr)
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil(), fmt.Errorf("invalid %s: %v: %w", param, err, domainErr)
	}

	if id == uuid.Nil() {
		return uuid.Nil(), fmt.Errorf("%s must not be nil UUID: %w", param, domainErr)
	}

	return id, nil
}

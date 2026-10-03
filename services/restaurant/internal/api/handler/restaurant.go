package handler

import (
	"context"
	"net/http"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/middleware"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"go.uber.org/zap"
)

//go:generate go tool minimock -i RestaurantService -o restaurant_service_mock_test.go
type RestaurantService interface {
	PutRestaurant(context.Context, domain.Restaurant) (domain.Restaurant, error)
	GetRestaurants(context.Context) ([]domain.Restaurant, error)
	GetRestaurant(context.Context, uuid.UUID) (domain.Restaurant, error)
	PutMenuItem(context.Context, uuid.UUID, domain.MenuItem) (domain.MenuItem, error)
	GetItems(context.Context, uuid.UUID) ([]domain.MenuItem, error)
	GetItem(context.Context, uuid.UUID, uuid.UUID) (domain.MenuItem, error)
}

type RestaurantHTTPHandler struct {
	restaurantService RestaurantService
	//TODO(review): теперь у этой структуры поле с логером
	// в каждом хэндлере я донастраиваю логер, но везде донастраиваю его одинаково.
	// кажется класть логер в контекст было не таким плохим решением??
	// решил через хэлпер requestLogger, норм?
	logger *zap.Logger
}

func NewRestaurantHTTPHandler(
	restaurantService RestaurantService,
	logger *zap.Logger,
) *RestaurantHTTPHandler {
	return &RestaurantHTTPHandler{
		restaurantService: restaurantService,
		logger:            logger,
	}
}

func (h *RestaurantHTTPHandler) requestLogger(r *http.Request) *zap.Logger {
	return h.logger.With(
		zap.String("request_id", r.Header.Get("X-Request-ID")),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)
}

func (h *RestaurantHTTPHandler) Routes() []api.Route {
	return []api.Route{
		{
			Path:    "/restaurants/{restaurant_id}",
			Method:  http.MethodPut,
			Handler: h.PutRestaurant,
			Middleware: []middleware.Middleware{
				middleware.RestaurantID(h.logger),
			},
		},
		{
			Path:    "/restaurants",
			Method:  http.MethodGet,
			Handler: h.GetRestaurants,
		},
		{
			Path:    "/restaurants/{restaurant_id}",
			Method:  http.MethodGet,
			Handler: h.GetRestaurant,
		},
		{
			Path:    "/restaurants/{restaurant_id}/menu/{menu_item_id}",
			Method:  http.MethodPut,
			Handler: h.PutMenuItem,
			Middleware: []middleware.Middleware{
				middleware.RestaurantID(h.logger),
			},
		},
		{
			Path:    "/restaurants/{restaurant_id}/menu",
			Method:  http.MethodGet,
			Handler: h.GetItems,
		},
		{
			Path:    "/restaurants/{restaurant_id}/menu/{menu_item_id}",
			Method:  http.MethodGet,
			Handler: h.GetItem,
		},
	}
}

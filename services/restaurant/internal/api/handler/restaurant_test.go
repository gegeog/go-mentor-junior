package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRestaurantHTTPHandler_GetRestaurant_OK(t *testing.T) {
	restaurantID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	restaurant := domain.NewRestaurant(restaurantID, "kfc", true, 1000, "USD")

	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantMock.
		Expect(minimock.AnyContext, restaurantID).
		Return(restaurant, nil)

	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+restaurantID.String(), nil)
	req.SetPathValue("restaurant_id", restaurantID.String())

	recorder := httptest.NewRecorder()
	handler.GetRestaurant(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
			"id": "11111111-1111-4111-8111-111111111111",
			"name": "kfc",
			"accepting_orders": true,
			"minimum_order_minor": 1000,
			"currency": "USD"
		}`, recorder.Body.String())
}

func TestKubernetesStubRoutes(t *testing.T) {
	router := api.NewRouter()
	router.RegisterRoutes(GetKuberStubRoutes(fakePinger{})...)

	paths := []string{"/livez"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Empty(t, recorder.Body.Bytes())
		})
	}
}

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(context.Context) error {
	return f.err
}

func TestReadyz_PingError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	p := fakePinger{err: errors.New("err")}
	readyzHandler(p)(rec, req)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestReadyz_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	p := fakePinger{err: nil}
	readyzHandler(p)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestRestaurantHTTPHandler_GetRestaurant_InvalidPath(t *testing.T) {
	svc := NewRestaurantServiceMock(t)
	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())

	req := httptest.NewRequest(http.MethodGet, "/restaurants/123", nil)
	req.SetPathValue("restaurant_id", "123")

	recorder := httptest.NewRecorder()

	handler.GetRestaurant(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	var body response.ErrorResponseDTO
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(t, "INVALID_RESTAURANT", body.Error.Code)
}

func TestRestaurantHTTPHandler_GetRestaurant_NotFound(t *testing.T) {
	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantMock.Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)
	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())

	restaurantID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+restaurantID.String(), nil)
	req.SetPathValue("restaurant_id", restaurantID.String())

	recorder := httptest.NewRecorder()

	handler.GetRestaurant(recorder, req)

	require.Equal(t, http.StatusNotFound, recorder.Code)
	var body response.ErrorResponseDTO
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(t, "RESTAURANT_NOT_FOUND", body.Error.Code)
}

func TestRestaurantHTTPHandler_GetRestaurants_Empty(t *testing.T) {
	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantsMock.Return([]domain.Restaurant{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/restaurants", nil)
	recorder := httptest.NewRecorder()

	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
	handler.GetRestaurants(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var body RestaurantsResponseDTO
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Empty(t, body.Restaurants)
}

func TestRestaurantHTTPHandler_GetRestaurants_OK(t *testing.T) {
	kfc := domain.NewRestaurant(uuid.New(), "kfc", true, 1000, "USD")
	mcdonalds := domain.NewRestaurant(uuid.New(), "mcd", false, 500, "RUB")

	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantsMock.Return([]domain.Restaurant{kfc, mcdonalds}, nil)

	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/restaurants", nil)
	recorder := httptest.NewRecorder()

	handler.GetRestaurants(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"restaurants": [
			{
				"id": "`+kfc.ID.String()+`",
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": 1000,
				"currency": "USD"
			},
			{
				"id": "`+mcdonalds.ID.String()+`",
				"name": "mcd",
				"accepting_orders": false,
				"minimum_order_minor": 500,
				"currency": "RUB"
			}
		]
	}`, recorder.Body.String())
}

func TestRestaurantHTTPHandler_GetItem_OK(t *testing.T) {
	restaurantID := uuid.New()
	itemID := uuid.New()
	item := domain.NewMenuItem(itemID, "wings", "spicy", 500, "USD", true)

	svc := NewRestaurantServiceMock(t)
	svc.GetItemMock.
		Expect(minimock.AnyContext, restaurantID, itemID).
		Return(item, nil)

	reqPath := fmt.Sprintf("/restaurants/%s/menu/%s", restaurantID.String(), itemID.String())
	req := httptest.NewRequest(http.MethodGet, reqPath, nil)
	req.SetPathValue("restaurant_id", restaurantID.String())
	req.SetPathValue("menu_item_id", itemID.String())

	recorder := httptest.NewRecorder()

	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
	handler.GetItem(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
    "id": "`+itemID.String()+`",
    "name": "wings",
    "description": "spicy",
    "price_minor": 500,
    "currency": "USD",
    "available": true
	}`, recorder.Body.String())
}

func TestRestaurantHTTPHandler_GetItem_Errors(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*RestaurantServiceMock)
		req        func() *http.Request
		wantStatus int
		wantCode   string
	}{
		{
			name:  "invalid restaurant uuid",
			setup: func(rsm *RestaurantServiceMock) {},
			req: func() *http.Request {
				itemID := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/123/menu/"+itemID.String(),
					nil,
				)
				req.SetPathValue("restaurant_id", "123")
				req.SetPathValue("menu_item_id", itemID.String())

				return req
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:  "invalid menu item uuid",
			setup: func(rsm *RestaurantServiceMock) {},
			req: func() *http.Request {
				restaurantID := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/"+restaurantID.String()+"/menu/123",
					nil,
				)
				req.SetPathValue("restaurant_id", restaurantID.String())
				req.SetPathValue("menu_item_id", "123")

				return req
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name: "restaurant not found",
			setup: func(rsm *RestaurantServiceMock) {
				rsm.GetItemMock.Return(domain.MenuItem{}, domain.ErrRestaurantNotFound)
			},
			req: func() *http.Request {
				restaurantID := uuid.New()
				itemID := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/"+restaurantID.String()+"/menu/"+itemID.String(),
					nil,
				)
				req.SetPathValue("restaurant_id", restaurantID.String())
				req.SetPathValue("menu_item_id", itemID.String())

				return req
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "RESTAURANT_NOT_FOUND",
		},
		{
			name: "menu item not found",
			setup: func(rsm *RestaurantServiceMock) {
				rsm.GetItemMock.Return(domain.MenuItem{}, domain.ErrMenuItemNotFound)
			},
			req: func() *http.Request {
				restaurantID := uuid.New()
				itemID := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/"+restaurantID.String()+"/menu/"+itemID.String(),
					nil,
				)
				req.SetPathValue("restaurant_id", restaurantID.String())
				req.SetPathValue("menu_item_id", itemID.String())

				return req
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "MENU_ITEM_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewRestaurantServiceMock(t)
			tt.setup(svc)

			req := tt.req()
			recorder := httptest.NewRecorder()

			handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
			handler.GetItem(recorder, req)

			require.Equal(t, tt.wantStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestRestaurantHTTPHandler_GetItems_OK(t *testing.T) {
	wingsID := uuid.New()
	wings := domain.NewMenuItem(wingsID, "wings", "spicy", 500, "USD", true)
	sauceID := uuid.New()
	sauce := domain.NewMenuItem(sauceID, "sauce", "cheese", 100, "USD", true)
	tests := []struct {
		name       string
		req        func() (*http.Request, uuid.UUID)
		setup      func(m *RestaurantServiceMock)
		wantStatus int
		wantBody   func(rid uuid.UUID) string
	}{
		{
			name: "not empty list",
			req: func() (*http.Request, uuid.UUID) {
				restaurantID := uuid.New()
				req := httptest.NewRequest(http.MethodGet, "/restaurants/"+restaurantID.String()+"/menu", nil)
				req.SetPathValue("restaurant_id", restaurantID.String())
				return req, restaurantID
			},
			setup: func(m *RestaurantServiceMock) {
				m.GetItemsMock.Return([]domain.MenuItem{wings, sauce}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(restaurantID uuid.UUID) string {
				return `{
					"restaurant_id": "` + restaurantID.String() + `",
					"items": [
						{
							"id": "` + wingsID.String() + `",
							"name": "wings",
							"description": "spicy",
							"price_minor": 500,
							"currency": "USD",
							"available": true
						},
						{
							"id": "` + sauceID.String() + `",
							"name": "sauce",
							"description": "cheese",
							"price_minor": 100,
							"currency": "USD",
							"available": true
						}
					]
				}`
			},
		},
		{
			name: "empty list",
			req: func() (*http.Request, uuid.UUID) {
				restaurantID := uuid.New()
				req := httptest.NewRequest(http.MethodGet, "/restaurants/"+restaurantID.String()+"/menu", nil)
				req.SetPathValue("restaurant_id", restaurantID.String())
				return req, restaurantID
			},
			setup: func(m *RestaurantServiceMock) {
				m.GetItemsMock.Return([]domain.MenuItem{}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(restaurantID uuid.UUID) string {
				return `{
					"restaurant_id": "` + restaurantID.String() + `",
					"items": []
				}`
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewRestaurantServiceMock(t)
			tt.setup(svc)
			req, rid := tt.req()
			recorder := httptest.NewRecorder()
			handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
			handler.GetItems(recorder, req)
			require.Equal(t, tt.wantStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody(rid), recorder.Body.String())
		})
	}
}

func TestRestaurantHTTPHandler_GetItems_Error(t *testing.T) {
	tests := []struct {
		name       string
		req        func() *http.Request
		setup      func(m *RestaurantServiceMock)
		wantStatus int
		wantCode   string
	}{
		{
			name: "invalid restaurant uuid",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/restaurants/123/menu", nil)
				req.SetPathValue("restaurant_id", "123")
				return req
			},
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name: "restaurant not found",
			req: func() *http.Request {
				restaurantID := uuid.New()
				req := httptest.NewRequest(http.MethodGet, "/restaurants/"+restaurantID.String()+"/menu", nil)
				req.SetPathValue("restaurant_id", restaurantID.String())
				return req
			},
			setup: func(m *RestaurantServiceMock) {
				m.GetItemsMock.Return([]domain.MenuItem{}, domain.ErrRestaurantNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "RESTAURANT_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewRestaurantServiceMock(t)
			tt.setup(svc)

			h := NewRestaurantHTTPHandler(svc, zap.NewNop())

			req := tt.req()
			recorder := httptest.NewRecorder()

			h.GetItems(recorder, req)

			require.Equal(t, tt.wantStatus, recorder.Code)
			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestRestaurantHTTPHandler_PutRestaurant_Errors(t *testing.T) {
	validBody := `{
		"name": "kfc",
		"accepting_orders": true,
		"minimum_order_minor": 1000,
		"currency": "USD"
	}`

	tests := []struct {
		name       string
		pathID     string
		reqBody    string
		setup      func(m *RestaurantServiceMock)
		wantStatus int
		wantCode   string
	}{
		{
			name:       "invalid restaurant uuid",
			pathID:     "not-a-uuid",
			reqBody:    validBody,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "empty body",
			reqBody:    ``,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "invalid json",
			reqBody:    `{"name":"kfc",,,}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name: "extra field in body",
			reqBody: `{
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": 1000,
				"currency": "USD",
				"unknown_field": 123
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name: "two jsons in body",
			reqBody: `{
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": 1000,
				"currency": "USD"
			}{
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": 1000,
				"currency": "USD"
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name: "missing required field",
			reqBody: `{
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": 1000
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name: "empty name",
			reqBody: `{
				"name": "",
				"accepting_orders": true,
				"minimum_order_minor": 1000,
				"currency": "USD"
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name: "negative minimum order minor",
			reqBody: `{
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": -12,
				"currency": "USD"
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:    "currency mismatch",
			reqBody: validBody,
			setup: func(m *RestaurantServiceMock) {
				m.PutRestaurantMock.Return(domain.Restaurant{}, domain.ErrCurrencyMismatch)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "CURRENCY_MISMATCH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewRestaurantServiceMock(t)
			tt.setup(svc)

			pathID := tt.pathID
			if pathID == "" {
				pathID = uuid.New().String()
			}

			req := httptest.NewRequest(http.MethodPut, "/restaurants/"+pathID, strings.NewReader(tt.reqBody))
			req.SetPathValue("restaurant_id", pathID)
			recorder := httptest.NewRecorder()

			handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
			handler.PutRestaurant(recorder, req)

			require.Equal(t, tt.wantStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestRestaurantHTTPHandler_PutRestaurant_OK(t *testing.T) {
	body := `{
    	"name":"kfc",
    	"accepting_orders":true,
    	"minimum_order_minor":1000,
    	"currency":"USD"
	}`

	restaurantID := uuid.New()
	restaurant := domain.NewRestaurant(restaurantID, "kfc", true, 1000, "USD")

	svc := NewRestaurantServiceMock(t)
	svc.PutRestaurantMock.Return(restaurant, nil)

	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+restaurantID.String(), strings.NewReader(body))
	req.SetPathValue("restaurant_id", restaurantID.String())
	recorder := httptest.NewRecorder()

	handler.PutRestaurant(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
			"id": "`+restaurantID.String()+`",
			"name": "kfc",
			"accepting_orders": true,
			"minimum_order_minor":1000,
			"currency":"USD"
		}`, recorder.Body.String())
}

func TestRestaurantHTTPHandler_PutMenuItem_Errors(t *testing.T) {
	tests := []struct {
		name             string
		restaurantPathID string
		itemPathID       string
		reqBody          string
		setup            func(m *RestaurantServiceMock)
		wantStatus       int
		wantCode         string
	}{
		{
			name:             "invalid restaurant uuid",
			restaurantPathID: "not-a-uuid",
			itemPathID:       "uuid",
			setup:            func(m *RestaurantServiceMock) {},
			wantStatus:       http.StatusBadRequest,
			wantCode:         "INVALID_RESTAURANT",
		},
		{
			name:             "invalid item uuid",
			restaurantPathID: uuid.New().String(),
			itemPathID:       "not-a-uuid",
			setup:            func(m *RestaurantServiceMock) {},
			wantStatus:       http.StatusBadRequest,
			wantCode:         "INVALID_MENU_ITEM",
		},
		{
			name:             "empty name",
			restaurantPathID: uuid.New().String(),
			itemPathID:       uuid.New().String(),
			reqBody: `
			{
			    "name":"",
			    "description":"big",
			    "price_minor":1116,
			    "currency":"usdt",
			    "available":false
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name:             "missing field",
			restaurantPathID: uuid.New().String(),
			itemPathID:       uuid.New().String(),
			reqBody: `
			{
			   "name":"shefroll",
			   "price_minor":1116,
			   "currency":"usdt",
			   "available":false
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name:             "negative proce minor",
			restaurantPathID: uuid.New().String(),
			itemPathID:       uuid.New().String(),
			reqBody: `
			{
			    "name":"shefroll",
			    "description":"big",
			    "price_minor":-1116,
			    "currency":"usdt",
			    "available":false
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name:             "two jsons",
			restaurantPathID: uuid.New().String(),
			itemPathID:       uuid.New().String(),
			reqBody: `
			{
			    "name":"shefroll",
			    "description":"big",
			    "price_minor":116,
			    "currency":"usdt",
			    "available":false
			}{
			    "name":"shefroll",
			    "description":"big",
			    "price_minor":116,
			    "currency":"usdt",
			    "available":false
			}`,
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name:             "restaurant not found",
			restaurantPathID: uuid.New().String(),
			itemPathID:       uuid.New().String(),
			reqBody: `
			{
			    "name":"shefroll",
			    "description":"big",
			    "price_minor":116,
			    "currency":"usdt",
			    "available":false
			}`,
			setup: func(m *RestaurantServiceMock) {
				m.PutMenuItemMock.Return(domain.MenuItem{}, domain.ErrRestaurantNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantCode:   "RESTAURANT_NOT_FOUND",
		},
		{
			name:             "currency mismatch",
			restaurantPathID: uuid.New().String(),
			itemPathID:       uuid.New().String(),
			reqBody: `
			{
			    "name":"shefroll",
			    "description":"big",
			    "price_minor":116,
			    "currency":"usdt",
			    "available":false
			}`,
			setup: func(m *RestaurantServiceMock) {
				m.PutMenuItemMock.Return(domain.MenuItem{}, domain.ErrCurrencyMismatch)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "CURRENCY_MISMATCH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewRestaurantServiceMock(t)
			tt.setup(svc)

			req := httptest.NewRequest(http.MethodPut, "/restaurants/"+tt.restaurantPathID+"/menu/"+tt.itemPathID, strings.NewReader(tt.reqBody))
			req.SetPathValue("restaurant_id", tt.restaurantPathID)
			req.SetPathValue("menu_item_id", tt.itemPathID)

			recorder := httptest.NewRecorder()

			handler := NewRestaurantHTTPHandler(svc, zap.NewNop())
			handler.PutMenuItem(recorder, req)

			require.Equal(t, tt.wantStatus, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestRestaurantHTTPHandler_PutMenuItem_OK(t *testing.T) {
	body := `{
	    "name":"shefroll",
	    "description":"очень вкусно",
	    "price_minor":66,
	    "currency":"usdt",
	    "available":false
	}`

	itemID := uuid.New()
	item := domain.NewMenuItem(itemID, "shefroll", "очень вкусно", 66, "usdt", false)

	restaurantID := uuid.New()
	svc := NewRestaurantServiceMock(t)
	svc.PutMenuItemMock.
		Expect(minimock.AnyContext, restaurantID, item).
		Return(item, nil)

	handler := NewRestaurantHTTPHandler(svc, zap.NewNop())

	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+restaurantID.String()+"/menu/"+itemID.String(), strings.NewReader(body))
	req.SetPathValue("restaurant_id", restaurantID.String())
	req.SetPathValue("menu_item_id", itemID.String())

	recorder := httptest.NewRecorder()

	handler.PutMenuItem(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
			"id": "`+itemID.String()+`",
			"name": "shefroll",
			"description":"очень вкусно",
			"price_minor":66,
		    "currency":"usdt",
		    "available":false
		}`, recorder.Body.String())
}

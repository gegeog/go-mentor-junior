package handler

import (
	"encoding/json"
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

func TestGetRestaurant_OK(t *testing.T) {
	rid := uuid.New()
	restaurant := domain.NewRestaurant(rid, "kfc", true, 1000, "USD")
	svc := NewRestaurantServiceMock(t)

	svc.GetRestaurantMock.
		Expect(minimock.AnyContext, rid).
		Return(restaurant, nil)

	h := NewRestaurantHTTPHandler(svc, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+rid.String(), nil)
	req.SetPathValue("restaurant_id", rid.String())

	rec := httptest.NewRecorder()

	h.GetRestaurant(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
			"id": "`+rid.String()+`",
			"name": "kfc",
			"accepting_orders": true,
			"minimum_order_minor": 1000,
			"currency": "USD"
		}`, rec.Body.String())
}

func TestKubernetesStubRoutes(t *testing.T) {
	router := api.NewRouter()
	router.RegisterRoutes(GetKuberStubRoutes()...)

	paths := []string{"/livez", "/readyz"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Empty(t, rec.Body.Bytes())
		})
	}
}

func TestGetRestaurant_InvalidPath(t *testing.T) {
	svc := NewRestaurantServiceMock(t)
	h := NewRestaurantHTTPHandler(svc, zap.NewNop())

	req := httptest.NewRequest(http.MethodGet, "/restaurants/123", nil)
	req.SetPathValue("restaurant_id", "123")

	rec := httptest.NewRecorder()

	h.GetRestaurant(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var body response.ErrorResponseDTO
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "INVALID_RESTAURANT", body.Error.Code)
}

func TestGetRestaurant_NotFound(t *testing.T) {
	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantMock.Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)
	h := NewRestaurantHTTPHandler(svc, zap.NewNop())

	rid := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+rid.String(), nil)
	req.SetPathValue("restaurant_id", rid.String())

	rec := httptest.NewRecorder()

	h.GetRestaurant(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	var body response.ErrorResponseDTO
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "RESTAURANT_NOT_FOUND", body.Error.Code)
}

func TestGetRestaurants_Empty(t *testing.T) {
	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantsMock.Return([]domain.Restaurant{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/restaurants", nil)
	rec := httptest.NewRecorder()

	h := NewRestaurantHTTPHandler(svc, zap.NewNop())
	h.GetRestaurants(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body RestaurantsResponseDTO
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Empty(t, body.Restaurants)
}

func TestGetRestaurants_OK(t *testing.T) {
	r1 := domain.NewRestaurant(uuid.New(), "kfc", true, 1000, "USD")
	r2 := domain.NewRestaurant(uuid.New(), "mcd", false, 500, "RUB")

	svc := NewRestaurantServiceMock(t)
	svc.GetRestaurantsMock.Return([]domain.Restaurant{r1, r2}, nil)

	h := NewRestaurantHTTPHandler(svc, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/restaurants", nil)
	rec := httptest.NewRecorder()

	h.GetRestaurants(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"restaurants": [
			{
				"id": "`+r1.ID.String()+`",
				"name": "kfc",
				"accepting_orders": true,
				"minimum_order_minor": 1000,
				"currency": "USD"
			},
			{
				"id": "`+r2.ID.String()+`",
				"name": "mcd",
				"accepting_orders": false,
				"minimum_order_minor": 500,
				"currency": "RUB"
			}
		]
	}`, rec.Body.String())
}

func TestGetItem_OK(t *testing.T) {
	rid := uuid.New()
	iid := uuid.New()
	item := domain.NewMenuItem(iid, "wings", "spicy", 500, "USD", true)

	svc := NewRestaurantServiceMock(t)

	svc.GetItemMock.
		Expect(minimock.AnyContext, rid, iid).
		Return(item, nil)

	reqPath := fmt.Sprintf("/restaurants/%s/menu/%s", rid.String(), iid.String())
	req := httptest.NewRequest(http.MethodGet, reqPath, nil)
	req.SetPathValue("restaurant_id", rid.String())
	req.SetPathValue("menu_item_id", iid.String())

	rec := httptest.NewRecorder()

	h := NewRestaurantHTTPHandler(svc, zap.NewNop())
	h.GetItem(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
    "id": "`+iid.String()+`",
    "name": "wings",
    "description": "spicy",
    "price_minor": 500,
    "currency": "USD",
    "available": true
	}`, rec.Body.String())
}

func TestGetItem_Errors(t *testing.T) {
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
				id := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/123/menu/"+id.String(),
					nil,
				)
				req.SetPathValue("restaurant_id", "123")
				req.SetPathValue("menu_item_id", id.String())

				return req
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:  "invalid menu item uuid",
			setup: func(rsm *RestaurantServiceMock) {},
			req: func() *http.Request {
				id := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/"+id.String()+"/menu/123",
					nil,
				)
				req.SetPathValue("restaurant_id", id.String())
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
				rid := uuid.New()
				iid := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/"+rid.String()+"/menu/"+iid.String(),
					nil,
				)
				req.SetPathValue("restaurant_id", rid.String())
				req.SetPathValue("menu_item_id", iid.String())

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
				rid := uuid.New()
				iid := uuid.New()
				req := httptest.NewRequest(
					http.MethodGet,
					"/restaurants/"+rid.String()+"/menu/"+iid.String(),
					nil,
				)
				req.SetPathValue("restaurant_id", rid.String())
				req.SetPathValue("menu_item_id", iid.String())

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
			rec := httptest.NewRecorder()

			h := NewRestaurantHTTPHandler(svc, zap.NewNop())
			h.GetItem(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestGetItems_OK(t *testing.T) {
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
				rid := uuid.New()
				req := httptest.NewRequest(http.MethodGet, "/restaurants/"+rid.String()+"/menu", nil)
				req.SetPathValue("restaurant_id", rid.String())
				return req, rid
			},
			setup: func(m *RestaurantServiceMock) {
				m.GetItemsMock.Return([]domain.MenuItem{wings, sauce}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(rid uuid.UUID) string {
				return `{
					"restaurant_id": "` + rid.String() + `",
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
				rid := uuid.New()
				req := httptest.NewRequest(http.MethodGet, "/restaurants/"+rid.String()+"/menu", nil)
				req.SetPathValue("restaurant_id", rid.String())
				return req, rid
			},
			setup: func(m *RestaurantServiceMock) {
				m.GetItemsMock.Return([]domain.MenuItem{}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(rid uuid.UUID) string {
				return `{
					"restaurant_id": "` + rid.String() + `",
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
			rec := httptest.NewRecorder()
			h := NewRestaurantHTTPHandler(svc, zap.NewNop())
			h.GetItems(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody(rid), rec.Body.String())
		})
	}
}

func TestGetItems_Error(t *testing.T) {
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
				rid := uuid.New()
				req := httptest.NewRequest(http.MethodGet, "/restaurants/"+rid.String()+"/menu", nil)
				req.SetPathValue("restaurant_id", rid.String())
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
			rec := httptest.NewRecorder()

			h.GetItems(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestPutRestaurant_Errors(t *testing.T) {
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
			rec := httptest.NewRecorder()

			h := NewRestaurantHTTPHandler(svc, zap.NewNop())
			h.PutRestaurant(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestPutRestaurant_OK(t *testing.T) {
	body := `{
    	"name":"kfc",
    	"accepting_orders":true,
    	"minimum_order_minor":1000,
    	"currency":"USD"
	}`

	rid := uuid.New()
	restaurant := domain.NewRestaurant(rid, "kfc", true, 1000, "USD")

	svc := NewRestaurantServiceMock(t)
	svc.PutRestaurantMock.Return(restaurant, nil)

	h := NewRestaurantHTTPHandler(svc, zap.NewNop())
	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+rid.String(), strings.NewReader(body))
	req.SetPathValue("restaurant_id", rid.String())
	rec := httptest.NewRecorder()

	h.PutRestaurant(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
			"id": "`+rid.String()+`",
			"name": "kfc",
			"accepting_orders": true,
			"minimum_order_minor":1000,
			"currency":"USD"
		}`, rec.Body.String())
}

func TestPutItem_Errors(t *testing.T) {
	tests := []struct {
		name       string
		rPathID    string
		iPathID    string
		reqBody    string
		setup      func(m *RestaurantServiceMock)
		wantStatus int
		wantCode   string
	}{
		{
			name:       "invalid restaurant uuid",
			rPathID:    "not-a-uuid",
			iPathID:    "uuid",
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "invalid item uuid",
			rPathID:    uuid.New().String(),
			iPathID:    "not-a-uuid",
			setup:      func(m *RestaurantServiceMock) {},
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name:    "empty name",
			rPathID: uuid.New().String(),
			iPathID: uuid.New().String(),
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
			name:    "missing field",
			rPathID: uuid.New().String(),
			iPathID: uuid.New().String(),
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
			name:    "negative proce minor",
			rPathID: uuid.New().String(),
			iPathID: uuid.New().String(),
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
			name:    "two jsons",
			rPathID: uuid.New().String(),
			iPathID: uuid.New().String(),
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
			name:    "restaurant not found",
			rPathID: uuid.New().String(),
			iPathID: uuid.New().String(),
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
			name:    "currency mismatch",
			rPathID: uuid.New().String(),
			iPathID: uuid.New().String(),
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

			req := httptest.NewRequest(http.MethodPut, "/restaurants/"+tt.rPathID+"/menu/"+tt.iPathID, strings.NewReader(tt.reqBody))
			req.SetPathValue("restaurant_id", tt.rPathID)
			req.SetPathValue("menu_item_id", tt.iPathID)

			rec := httptest.NewRecorder()

			h := NewRestaurantHTTPHandler(svc, zap.NewNop())
			h.PutMenuItem(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body response.ErrorResponseDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
		})
	}
}

func TestPutMenuItem_OK(t *testing.T) {
	body := `{
	    "name":"shefroll",
	    "description":"очень вкусно",
	    "price_minor":66,
	    "currency":"usdt",
	    "available":false
	}`

	iid := uuid.New()
	item := domain.NewMenuItem(iid, "shefroll", "очень вкусно", 66, "usdt", false)

	rid := uuid.New()
	svc := NewRestaurantServiceMock(t)
	svc.PutMenuItemMock.
		Expect(minimock.AnyContext, rid, item).
		Return(item, nil)

	h := NewRestaurantHTTPHandler(svc, zap.NewNop())

	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+rid.String()+"/menu/"+iid.String(), strings.NewReader(body))
	req.SetPathValue("restaurant_id", rid.String())
	req.SetPathValue("menu_item_id", iid.String())

	rec := httptest.NewRecorder()

	h.PutMenuItem(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
			"id": "`+iid.String()+`",
			"name": "shefroll",
			"description":"очень вкусно",
			"price_minor":66,
		    "currency":"usdt",
		    "available":false
		}`, rec.Body.String())
}

package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "invalid menu item",
			err:        domain.ErrInvalidMenuItem,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_MENU_ITEM",
		},
		{
			name:       "invalid restaurant",
			err:        domain.ErrInvalidRestaurant,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "restaurant access denied",
			err:        domain.ErrRestaurantAccessDenied,
			wantStatus: http.StatusForbidden,
			wantCode:   "RESTAURANT_ACCESS_DENIED",
		},
		{
			name:       "restaurant not found",
			err:        domain.ErrRestaurantNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   "RESTAURANT_NOT_FOUND",
		},
		{
			name:       "menu item not found",
			err:        domain.ErrMenuItemNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   "MENU_ITEM_NOT_FOUND",
		},
		{
			name:       "currency mismatch",
			err:        domain.ErrCurrencyMismatch,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "CURRENCY_MISMATCH",
		},
		{
			name:       "wrapped error",
			err:        fmt.Errorf("put menu item: %w", domain.ErrCurrencyMismatch),
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "CURRENCY_MISMATCH",
		},
		{
			name:       "unexpected error",
			err:        errors.New("database connection lost"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			ErrorResponse(
				zap.NewNop(),
				rec,
				tt.err,
				"test err msg",
			)

			require.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body ErrorResponseDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

			assert.Equal(t, tt.wantCode, body.Error.Code)
			assert.Equal(t, "test err msg", body.Error.Message)
			assert.Empty(t, body.Error.Details)
		})
	}
}

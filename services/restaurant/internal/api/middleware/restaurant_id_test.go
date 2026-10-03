package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRestaurantID(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name       string
		header     string
		pathID     string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "no header",
			header:     "",
			pathID:     validID,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "header uuid is nil",
			header:     "00000000-0000-0000-0000-000000000000",
			pathID:     validID,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "header not valid uuid",
			header:     "123",
			pathID:     validID,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "path not valid uuid",
			header:     validID,
			pathID:     "123",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_RESTAURANT",
		},
		{
			name:       "path not match with header",
			header:     uuid.New().String(),
			pathID:     validID,
			wantStatus: http.StatusForbidden,
			wantCode:   "RESTAURANT_ACCESS_DENIED",
		},
		{
			name:       "header match path",
			header:     validID,
			pathID:     validID,
			wantStatus: http.StatusOK,
			wantCode:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			h := RestaurantID(zap.NewNop())(handler)

			req := httptest.NewRequest(
				http.MethodPut,
				"/restaurants/"+tt.pathID,
				nil,
			)
			req.SetPathValue("restaurant_id", tt.pathID)
			if tt.header != "" {
				req.Header.Set("X-Restaurant-ID", tt.header)
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			callExpected := tt.wantCode == ""
			assert.Equal(t, callExpected, handlerCalled)

			if tt.wantCode != "" {
				var body response.ErrorResponseDTO
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
				assert.Equal(t, tt.wantCode, body.Error.Code)
			}
		})
	}
}

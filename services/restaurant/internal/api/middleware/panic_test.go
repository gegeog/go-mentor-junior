package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPanic(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
		wantCode   string
	}{
		{
			name: "no panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantStatus: http.StatusOK,
			wantCode:   "",
		},
		{
			name: "panic with 500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("AAAAAAAAAAAAAA!!!!")
			},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_SERVER_ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := Panic(zap.NewNop())
			handler := middleware(tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			require.Equal(t, tt.wantStatus, recorder.Code)
			if tt.wantCode != "" {
				var body response.ErrorResponseDTO
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
				assert.Equal(t, tt.wantCode, body.Error.Code)
			}
		})
	}
}

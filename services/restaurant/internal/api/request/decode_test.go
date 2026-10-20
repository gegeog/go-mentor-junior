package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type decodeTestRequest struct {
	Name string `json:"name" validate:"required"`
}

func TestDecodeAndValidate(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    decodeTestRequest
		wantErr bool
	}{
		{
			name: "valid request",
			body: `{"name":"kfc"}`,
			want: decodeTestRequest{
				Name: "kfc",
			},
		},
		{
			name:    "unknown field",
			body:    `{"name":"kfc","unknown":"value"}`,
			wantErr: true,
		},
		{
			name:    "two JSON objects",
			body:    `{"name":"kfc"}{"name":"bk"}`,
			wantErr: true,
		},
		{
			name:    "JSON followed by garbage",
			body:    `{"name":"kfc"} garbage`,
			wantErr: true,
		},
		{
			name:    "malformed JSON",
			body:    `{"name":`,
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    "",
			wantErr: true,
		},
		{
			name:    "validation error",
			body:    `{"name":""}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPut,
				"/restaurants/test",
				strings.NewReader(tt.body),
			)

			got, err := DecodeAndValidate[decodeTestRequest](req)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, domain.ErrInvalidArgument)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

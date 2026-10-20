package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"uuid"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePathIDs(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name       string
		parser     func(*http.Request) (uuid.UUID, error)
		paramName  string
		paramValue string
		setParam   bool
		wantID     uuid.UUID
		wantErr    error
	}{
		{
			name:       "valid restaurant ID",
			parser:     ParseRestaurantID,
			paramName:  PathParamRestaurantID,
			paramValue: validID.String(),
			setParam:   true,
			wantID:     validID,
		},
		{
			name:      "missing restaurant ID",
			parser:    ParseRestaurantID,
			paramName: PathParamRestaurantID,
			wantErr:   domain.ErrInvalidRestaurant,
		},
		{
			name:       "invalid restaurant ID",
			parser:     ParseRestaurantID,
			paramName:  PathParamRestaurantID,
			paramValue: "not-a-uuid",
			setParam:   true,
			wantErr:    domain.ErrInvalidRestaurant,
		},
		{
			name:       "nil restaurant ID",
			parser:     ParseRestaurantID,
			paramName:  PathParamRestaurantID,
			paramValue: uuid.Nil().String(),
			setParam:   true,
			wantErr:    domain.ErrInvalidRestaurant,
		},
		{
			name:       "valid menu item ID",
			parser:     ParseMenuItemID,
			paramName:  PathParamMenuItemID,
			paramValue: validID.String(),
			setParam:   true,
			wantID:     validID,
		},
		{
			name:      "missing menu item ID",
			parser:    ParseMenuItemID,
			paramName: PathParamMenuItemID,
			wantErr:   domain.ErrInvalidMenuItem,
		},
		{
			name:       "invalid menu item ID",
			parser:     ParseMenuItemID,
			paramName:  PathParamMenuItemID,
			paramValue: "not-a-uuid",
			setParam:   true,
			wantErr:    domain.ErrInvalidMenuItem,
		},
		{
			name:       "nil menu item ID",
			parser:     ParseMenuItemID,
			paramName:  PathParamMenuItemID,
			paramValue: uuid.Nil().String(),
			setParam:   true,
			wantErr:    domain.ErrInvalidMenuItem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.setParam {
				req.SetPathValue(tt.paramName, tt.paramValue)
			}

			gotID, err := tt.parser(req)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}

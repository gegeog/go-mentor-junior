package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			domain.ErrInvalidArgument,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"validate request: %v: %w",
			err,
			domain.ErrInvalidArgument,
		)
	}

	return nil
}

package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidate[T any](r *http.Request) (T, error) {
	var dest T

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dest); err != nil {
		return dest, fmt.Errorf(
			"decode request JSON: %w: %w",
			err,
			domain.ErrInvalidArgument,
		)
	}

	var extra any
	switch err := dec.Decode(&extra); {
	case errors.Is(err, io.EOF):
	case err != nil:
		return dest, fmt.Errorf(
			"decode trailing JSON: %w: %w",
			err,
			domain.ErrInvalidArgument,
		)
	default:
		return dest, fmt.Errorf(
			"request body must contain exactly one JSON value: %w",
			domain.ErrInvalidArgument,
		)
	}

	if err := requestValidator.Struct(&dest); err != nil {
		return dest, fmt.Errorf(
			"validate request: %w: %w",
			err,
			domain.ErrInvalidArgument,
		)
	}

	return dest, nil
}

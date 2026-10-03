package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"go.uber.org/zap"
)

type ErrorDetail struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type ErrorResponseDTO struct {
	Error ErrorDetail `json:"error"`
}

func JSONResponse(
	log *zap.Logger,
	rw http.ResponseWriter,
	responseBody any,
	statusCode int,
) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	if err := json.NewEncoder(rw).Encode(responseBody); err != nil {
		log.Error("write HTTP response", zap.Error(err))
	}
}

func ErrorResponse(
	log *zap.Logger,
	rw http.ResponseWriter,
	err error,
	msg string,
) {
	var (
		statusCode int
		code       string
	)

	logFunc := log.Warn

	switch {
	case errors.Is(err, domain.ErrInvalidMenuItem):
		code = "INVALID_MENU_ITEM"
		statusCode = http.StatusBadRequest
	case errors.Is(err, domain.ErrInvalidRestaurant):
		code = "INVALID_RESTAURANT"
		statusCode = http.StatusBadRequest
	case errors.Is(err, domain.ErrRestaurantAccessDenied):
		code = "RESTAURANT_ACCESS_DENIED"
		statusCode = http.StatusForbidden
	case errors.Is(err, domain.ErrRestaurantNotFound):
		code = "RESTAURANT_NOT_FOUND"
		statusCode = http.StatusNotFound
	case errors.Is(err, domain.ErrMenuItemNotFound):
		code = "MENU_ITEM_NOT_FOUND"
		statusCode = http.StatusNotFound
	case errors.Is(err, domain.ErrCurrencyMismatch):
		code = "CURRENCY_MISMATCH"
		statusCode = http.StatusUnprocessableEntity
	default:
		code = "INTERNAL_SERVER_ERROR"
		statusCode = http.StatusInternalServerError
		logFunc = log.Error
	}

	logFunc(msg, zap.Error(err))

	response := ErrorResponseDTO{
		Error: ErrorDetail{
			Code:    code,
			Message: msg,
			Details: make(map[string]any),
		},
	}

	JSONResponse(log, rw, response, statusCode)
}

func PanicResponse(
	log *zap.Logger,
	rw http.ResponseWriter,
	p any,
	msg string,
) {
	err := fmt.Errorf("unexpected panic: %v", p)
	ErrorResponse(log, rw, err, msg)
}

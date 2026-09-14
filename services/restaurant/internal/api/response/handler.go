package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/domain"
	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/logger"
)

type HTTPResponseHandler struct {
	log *logger.Logger
	rw  http.ResponseWriter
}

type ErrorDetail struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func NewHTTPResponseHandler(
	log *logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", slog.Any("error", err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(
	err error,
	msg string,
) {
	var (
		statusCode int
		logFunc    func(string, ...any)
		code       string
	)

	switch {
	case errors.Is(err, domain.ErrInvalidMenuItem):
		code = "INVALID_MENU_ITEM"
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, domain.ErrInvalidRestaurant):
		code = "INVALID_RESTAURANT"
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, domain.ErrRestaurantAccessDenied):
		code = "RESTAURANT_ACCESS_DENIED"
		statusCode = http.StatusForbidden
		logFunc = h.log.Warn
	case errors.Is(err, domain.ErrRestaurantNotFound):
		code = "RESTAURANT_NOT_FOUND"
		statusCode = http.StatusNotFound
		logFunc = h.log.Warn
	case errors.Is(err, domain.ErrMenuItemNotFound):
		code = "MENU_ITEM_NOT_FOUND"
		statusCode = http.StatusNotFound
		logFunc = h.log.Warn
	case errors.Is(err, domain.ErrCurrencyMismatch):
		code = "CURRENCY_MISMATCH"
		statusCode = http.StatusUnprocessableEntity
		logFunc = h.log.Warn
	default:
		code = "INTERNAL_SERVER_ERROR"
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, slog.Any("error", err))

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: msg,
			Details: make(map[string]any),
		},
	}

	h.JSONResponse(response, statusCode)
}

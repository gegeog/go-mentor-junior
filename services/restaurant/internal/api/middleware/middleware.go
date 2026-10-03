package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

var requestIDHeader = "X-Request-ID"

func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i, _ := range slices.Backward(m) {
		h = m[i](h)
	}

	return h
}

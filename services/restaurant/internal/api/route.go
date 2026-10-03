package api

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api/middleware"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc

	Middleware []middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	h := middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)

	return h
}

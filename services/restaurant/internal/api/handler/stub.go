package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api"
)

type Pinger interface {
	Ping(context.Context) error
}

func GetKuberStubRoutes(pinger Pinger) []api.Route {
	return []api.Route{
		{
			Path:    "/livez",
			Method:  http.MethodGet,
			Handler: StubOKResponse,
		},
		{
			Path:    "/readyz",
			Method:  http.MethodGet,
			Handler: readyzHandler(pinger),
		},
	}
}

func readyzHandler(p Pinger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()

		if err := p.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func StubOKResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

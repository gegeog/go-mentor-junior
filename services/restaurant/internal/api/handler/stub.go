package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/api"
)

func GetKuberStubRoutes() []api.Route {
	return []api.Route{
		{
			Path:    "/livez",
			Method:  http.MethodGet,
			Handler: StubOKResponse,
		},
		{
			Path:    "/readyz",
			Method:  http.MethodGet,
			Handler: StubOKResponse,
		},
	}
}

func StubOKResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

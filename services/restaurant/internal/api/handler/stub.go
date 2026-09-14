package handler

import (
	"net/http"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/app"
)

func GetKuberStubRoutes() []app.Route {
	return []app.Route{
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

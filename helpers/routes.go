package helpers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/project-jano/user-service.go/api/metric"
	"github.com/project-jano/user-service.go/logger"
	"github.com/project-jano/user-service.go/model"
)

func AppendRoutes(router *mux.Router, routes model.Routes, measureInPrometheus bool, traceCalls bool) {
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if traceCalls {
			handler = logger.Logger(handler, route.Name)
		}

		if measureInPrometheus {
			handler = metric.PrometheusMiddleware(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

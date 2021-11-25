package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	Authenticated bool
}

type Routes []Route

func (api *API) appendRoutes(router *mux.Router, routes Routes, measureInPrometheus bool) {
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if api.Configuration.TraceCallsEnabled {
			handler = Logger(handler, route.Name)
		}

		if measureInPrometheus {
			handler = PrometheusMiddleware(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

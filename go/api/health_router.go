package api

import (
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func (api *API) addHealthRouter() {

	const path = "/health"

	healthRouter := mux.NewRouter().PathPrefix(path).Subrouter().StrictSlash(true)

	routes := Routes{
		Route{
			"Liveness",
			strings.ToUpper("Get"),
			"/liveness",
			api.Liveness,
			false,
		},
		Route{
			"Readiness",
			strings.ToUpper("Get"),
			"/readiness",
			api.Liveness,
			false,
		},
	}

	api.appendRoutes(healthRouter, routes, false)

	api.Router.PathPrefix(path).Handler(negroni.New(
		negroni.Wrap(healthRouter),
	))
}

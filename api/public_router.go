package api

import (
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func (api *API) addSecurityRouter() {

	const version = "/v2"
	const path = "/public"

	securityRouter := mux.NewRouter().PathPrefix(version).PathPrefix(path).Subrouter().StrictSlash(true)

	routes := Routes{
		Route{
			"GetServiceCertificate",
			strings.ToUpper("Get"),
			"/certificate",
			api.GetServiceCertificate,
			false,
		},

		Route{
			"GetServiceInfo",
			strings.ToUpper("Get"),
			"/service",
			api.GetServiceInfo,
			false,
		},
	}

	api.appendRoutes(securityRouter, routes, true)

	api.Router.PathPrefix(version).PathPrefix(path).Handler(negroni.New(
		negroni.Wrap(securityRouter),
	))
}

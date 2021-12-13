package pub

import (
	"strings"

	"github.com/project-jano/user-service.go/app"
	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/model"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type API struct {
	serverCertificate string
	fingerprint       string
}

const (
	APIVersion = "/v2"
	APIPath = "/public"
)

func NewPublicAPI(router *mux.Router, configuration app.APIConfiguration, fingerprint string) {
	publicAPI := API{
		serverCertificate: configuration.CertificatePEM,
		fingerprint:       fingerprint,
	}
	publicAPI.addToRouter(router, configuration.TraceCallsEnabled)
}

func (public *API) addToRouter(router *mux.Router, traceCalls bool) {

	securityRouter := mux.NewRouter().PathPrefix(APIVersion).PathPrefix(APIPath).Subrouter().StrictSlash(true)

	routes := model.Routes{
		model.Route{
			Name:        "GetServiceCertificate",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/certificate",
			HandlerFunc: public.GetServiceCertificate,
		},

		model.Route{
			Name:        "GetServiceInfo",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/service",
			HandlerFunc: public.GetServiceInfo,
		},
	}

	helpers.AppendRoutes(securityRouter, routes, false, traceCalls)

	router.PathPrefix(APIVersion).PathPrefix(APIPath).Handler(negroni.New(
		negroni.Wrap(securityRouter),
	))
}

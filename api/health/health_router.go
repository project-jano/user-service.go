package health

import (
	"strings"

	"github.com/project-jano/user-service.go/app"
	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/model"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type API struct {
	MongoClient *mongo.Client
}

const (
	APIPath = "/health"
)

func NewHealthAPI(router *mux.Router, configuration app.APIConfiguration, mongoClient *mongo.Client) {
	healthAPI := API{
		MongoClient: mongoClient,
	}
	healthAPI.addToRouter(router, configuration.TraceCallsEnabled)
}

func (healthAPI *API) addToRouter(router *mux.Router, traceCalls bool) {

	healthRouter := mux.NewRouter().PathPrefix(APIPath).Subrouter().StrictSlash(true)

	routes := model.Routes{
		model.Route{
			Name:        "Liveness",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/liveness",
			HandlerFunc: healthAPI.Liveness,
		},
		model.Route{
			Name:        "Readiness",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/readiness",
			HandlerFunc: healthAPI.Liveness,
		},
	}

	helpers.AppendRoutes(healthRouter, routes, false, traceCalls)

	router.PathPrefix(APIPath).Handler(negroni.New(
		negroni.Wrap(healthRouter),
	))
}

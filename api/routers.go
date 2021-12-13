package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/project-jano/user-service.go/api/health"
	"github.com/project-jano/user-service.go/api/pub"
	"github.com/project-jano/user-service.go/api/user"
	"github.com/project-jano/user-service.go/app"
	"github.com/project-jano/user-service.go/logger"
	"github.com/project-jano/user-service.go/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	Configuration app.APIConfiguration
	Router        *mux.Router
	MongoClient   *mongo.Client
	UserDatabase  *mongo.Collection
	Fingerprint   string
}

func NewAPI(config app.APIConfiguration, client *mongo.Client, fingerprint string) *API {

	rootRouter := mux.NewRouter().StrictSlash(true)

	userDatabase := client.Database(config.DatabaseName).Collection("users")

	api := &API{
		Configuration: config,
		Router:        rootRouter,
		MongoClient:   client,
		UserDatabase:  userDatabase,
		Fingerprint:   fingerprint,
	}

	pub.NewPublicAPI(rootRouter, config, fingerprint)
	health.NewHealthAPI(rootRouter, config, client)
	user.NewUserAPI(rootRouter, config, fingerprint, repository.NewUserRepository(userDatabase))

	// Endpoints where prometheus Middleware doesnot apply
	api.Router.Handle("/metrics", promhttp.Handler())

	api.Router.NotFoundHandler = logger.Logger(http.NotFoundHandler(), "NotFound")

	return api
}

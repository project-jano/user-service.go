package health

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
	"context"
	"github.com/project-jano/user-service.go/helpers"
	"net/http"
	"os"
	"time"

	"github.com/project-jano/user-service.go/model"
)

const (
	ContentType        = "Content-Type"
	DefaultContentType = "application/json; charset=UTF-8"
)

func (healthAPI *API) Liveness(w http.ResponseWriter, _ *http.Request) {
	helpers.SetupDefaultContentType(w)

	hostname, ok := getHostname(w)
	if !ok {
		return
	}

	responseWithLivenessStatus(w, "up", hostname)
}

func (healthAPI *API) Readiness(w http.ResponseWriter, _ *http.Request) {
	helpers.SetupDefaultContentType(w)

	hostname, ok := getHostname(w)
	if !ok {
		return
	}

	// Check the connection to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBErr := healthAPI.MongoClient.Ping(ctx, nil)
	if mongoDBErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWithReadinessStatus(w, "ready", hostname)
}

func getHostname(w http.ResponseWriter) (string, bool) {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", false
	}
	return hostname, true
}

func responseWithLivenessStatus(w http.ResponseWriter, status string, hostname string) {
	liveness := model.LivenessResponse{
		Status:   status,
		Hostname: hostname,
	}

	helpers.ResponseWithJSON(w, liveness)
}

func responseWithReadinessStatus(w http.ResponseWriter, status string, hostname string) {
	readiness := model.ReadinessResponse{
		Status:   status,
		Hostname: hostname,
	}

	helpers.ResponseWithJSON(w, readiness)
}

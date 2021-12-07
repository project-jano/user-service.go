package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/project-jano/user-service.go/model"
)

func (a *API) Liveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(ContentType, DefaultContentType)

	hostname, ok := a.getHostname(w)
	if !ok {
		return
	}

	a.responseWithLivenessStatus(w, "up", hostname)
}

func (a *API) Readiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(ContentType, DefaultContentType)

	hostname, ok := a.getHostname(w)
	if !ok {
		return
	}

	// Check the connection to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBErr := a.MongoClient.Ping(ctx, nil)
	if mongoDBErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.responseWithReadinessStatus(w, "ready", hostname)
}

func (a *API) getHostname(w http.ResponseWriter) (string, bool) {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", false
	}
	return hostname, true
}

func (a *API) responseWithLivenessStatus(w http.ResponseWriter, status string, hostname string) {
	liveness := model.LivenessResponse{
		Status:   status,
		Hostname: hostname,
	}

	jsonErr := json.NewEncoder(w).Encode(liveness)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *API) responseWithReadinessStatus(w http.ResponseWriter, status string, hostname string) {
	readiness := model.ReadinessResponse{
		Status:   status,
		Hostname: hostname,
	}

	jsonErr := json.NewEncoder(w).Encode(readiness)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

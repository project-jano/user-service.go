package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func (a *API) Liveness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hostname, ok := a.getHostname(w)
	if !ok {
		return
	}

	a.responseWithStatus(w, "up", hostname)
}

func (a *API) Readiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

	a.responseWithStatus(w, "ready", hostname)
}

func (a *API) getHostname(w http.ResponseWriter) (string, bool) {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", false
	}
	return hostname, true
}

func (a *API) responseWithStatus(w http.ResponseWriter, status string, hostname string) {
	jsonErr := json.NewEncoder(w).Encode(map[string]string{"status": status, "hostname": hostname})
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

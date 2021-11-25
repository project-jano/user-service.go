package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/project-jano/user-service.go/go/model"
)

func (a *API) GetServiceCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	certificate := model.Certificate{
		Certificate: a.Configuration.CertificatePEM,
	}
	jsonErr := json.NewEncoder(w).Encode(certificate)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *API) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	serviceInfo := model.ServiceInfo{
		Hostname:    hostname,
		Timestamp:   time.Now().Unix(),
		Fingerprint: a.Fingerprint,
	}

	jsonErr := json.NewEncoder(w).Encode(serviceInfo)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

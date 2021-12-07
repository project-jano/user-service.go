package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/project-jano/user-service.go/model"
)

func (a *API) GetServiceCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, DefaultContentType)

	certificate := model.Certificate{
		Certificate: a.Configuration.CertificatePEM,
	}
	jsonErr := json.NewEncoder(w).Encode(certificate)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *API) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, DefaultContentType)

	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	serviceInfo := model.ServiceInformation{
		Hostname:    hostname,
		Timestamp:   time.Now().Unix(),
		Fingerprint: a.Fingerprint,
	}

	jsonErr := json.NewEncoder(w).Encode(serviceInfo)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

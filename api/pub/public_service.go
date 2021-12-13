package pub

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
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

func (public *API) GetServiceCertificate(w http.ResponseWriter, r *http.Request) {
	helpers.SetupDefaultContentType(w)

	certificate := model.Certificate{
		Certificate: public.serverCertificate,
	}

	helpers.ResponseWithJSON(w, certificate)
}

func (public *API) GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	helpers.SetupDefaultContentType(w)

	serviceInfo, ok := createServiceInfo(public.fingerprint)
	if !ok {
		return
	}

	helpers.ResponseWithJSON(w, serviceInfo)
}

func createServiceInfo(fingerprint string) (model.ServiceInformation, bool) {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		return model.ServiceInformation{}, false
	}
	serviceInfo := model.ServiceInformation{
		Hostname:    hostname,
		Timestamp:   time.Now().Unix(),
		Fingerprint: fingerprint,
	}
	return serviceInfo, true
}

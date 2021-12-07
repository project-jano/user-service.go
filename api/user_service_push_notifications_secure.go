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
)

func (a *API) SecurePushNotificationForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, DefaultContentType)

	params, paramsOk := a.getSecurePushNotificationRequestParams(w, r)
	if !paramsOk {
		return
	}

	user, userOk := a.findUser(params.userId, w)
	if !userOk {
		return
	}

	certificates := a.filterCertificates(user.Certificates, params.useAllDevices, params.useDefaultKey, params.devices, params.keyId)

	if len(certificates) == 0 {
		a.respondWithJSON(w, http.StatusBadRequest, "not found a valid device and key combination for this user")
		return
	}

	payload, jsonErr := json.Marshal(params.request)

	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	a.securePushNotification(w, r, string(payload), certificates, params.showCompleteOutput, params.showSplittedOutput)
}
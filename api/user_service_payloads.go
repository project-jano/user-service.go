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

	"github.com/project-jano/user-service.go/model"
)

func (a *API) DecodeSecurePayload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userId := extractUserId(r)
	if userId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid userId")
		return
	}

	deviceId := extractDeviceId(r)
	if deviceId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid deviceId")
		return
	}

	// Decode request
	var decodeSecurePayloadRequest model.DecodeSecurePayloadRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&decodeSecurePayloadRequest); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if !decodeSecurePayloadRequest.IsValid() {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

}

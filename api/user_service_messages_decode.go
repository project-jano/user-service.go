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

	"github.com/project-jano/user-service.go/model"
)

func (a *API) DecodeSecureMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, DefaultContentType)

	_, paramsError := a.getDecodeSecuredMessageRequestParams(w, r)

	if !paramsError {
		return
	}

}

func (a *API) getDecodeSecuredMessageRequestParams(w http.ResponseWriter, r *http.Request) (*decodeMessageRequestParams, bool) {

	userId := extractUserId(r)
	if userId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid userId")
		return nil, false
	}

	// Decode request
	var decodeSecurePayloadRequest model.DecodeSecuredMessageRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&decodeSecurePayloadRequest); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	if !decodeSecurePayloadRequest.IsValid() {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	useDefaultDevice := true
	deviceId := ""
	if deviceIdString := r.URL.Query().Get(QueryParamDeviceId); len(deviceIdString) > 0 {
		deviceId = deviceIdString
		useDefaultDevice = false
	}

	// Get keys in requests. If empty means default
	useDefaultKey := true
	keyId := DefaultKeyId
	if keyIdString := r.URL.Query().Get(QueryParamKeyId); len(keyIdString) > 0 {
		if keyIdString != DefaultKeyId {
			keyId = keyIdString
			useDefaultKey = false
		}
	}
	return &decodeMessageRequestParams{
		userId:           userId,
		deviceId:         deviceId,
		keyId:            keyId,
		useDefaultDevice: useDefaultDevice,
		useDefaultKey:    useDefaultKey,
		request:          decodeSecurePayloadRequest,
	}, true
}

type decodeMessageRequestParams struct {
	userId           string
	deviceId         string
	keyId            string
	useDefaultDevice bool
	useDefaultKey    bool
	request          model.DecodeSecuredMessageRequest
}

package user

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

	"github.com/project-jano/user-service.go/helpers"

	"github.com/project-jano/user-service.go/model"
)

func (userAPI *API) DecodeSecureMessage(w http.ResponseWriter, r *http.Request) {
	helpers.SetupDefaultContentType(w)

	//TODO
	_, _ = userAPI.getDecodeSecuredMessageRequestParams(w, r)

}

func (userAPI *API) getDecodeSecuredMessageRequestParams(w http.ResponseWriter, r *http.Request) (*decodeMessageRequestParams, bool) {

	// Decode request
	var decodeSecurePayloadRequest model.DecodeSecuredMessageRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&decodeSecurePayloadRequest); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	if !decodeSecurePayloadRequest.IsValid() {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	requestParams := extractUserParams(r, true)
	if requestParams == nil || len(requestParams.devices) != 1 {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid requests params")
	}

	return &decodeMessageRequestParams{
		userId:           requestParams.userId,
		deviceId:         requestParams.devices[0],
		keyId:            requestParams.keyId,
		useDefaultDevice: requestParams.useDefaultDevice,
		useDefaultKey:    requestParams.useDefaultKey,
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

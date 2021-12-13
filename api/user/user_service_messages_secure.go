package user

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
	"time"

	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/model"
)

func (userAPI *API) SecureMessageForUser(w http.ResponseWriter, r *http.Request) {
	helpers.SetupDefaultContentType(w)

	params, paramsOk := userAPI.getSecureMessageRequestParams(w, r)
	if !paramsOk {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := userAPI.repository.FindUser(ctx, params.userId)
	if user == nil {
		return
	}

	certificates := filterCertificates(user.Certificates, params.useAllDevices, params.useDefaultKey, params.devices, params.keyId)

	if len(certificates) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "not found a valid device and key combination for this user")
		return
	}

	userAPI.securePayload(w, params.request.Message, certificates)
}

func (userAPI *API) getSecureMessageRequestParams(w http.ResponseWriter, r *http.Request) (*secureMessageRequestParams, bool) {

	// Decode request
	var secureMessageRequest model.SecureMessageRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&secureMessageRequest); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	if !secureMessageRequest.IsValid() {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid secure message request")
		return nil, false
	}

	requestParams := extractUserParams(r, false)
	if requestParams == nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid requests params")
		return nil, false
	}

	return &secureMessageRequestParams{
		userId:           requestParams.userId,
		devices:          requestParams.devices,
		keyId:            requestParams.keyId,
		useAllDevices:    requestParams.useAllDevices,
		useDefaultDevice: requestParams.useDefaultDevice,
		useDefaultKey:    requestParams.useDefaultKey,
		request:          secureMessageRequest,
	}, true
}

type secureMessageRequestParams struct {
	userId           string
	devices          []string
	keyId            string
	useAllDevices    bool
	useDefaultDevice bool
	useDefaultKey    bool
	request          model.SecureMessageRequest
}

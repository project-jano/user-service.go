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

const (
	QueryParamShowCompleteOutput = "completeOutput"
	QueryParamShowSplittedOutput = "splittedOutput"
)

func (userAPI *API) SecurePushNotificationForUser(w http.ResponseWriter, r *http.Request) {
	helpers.SetupDefaultContentType(w)

	params, paramsOk := userAPI.getSecurePushNotificationRequestParams(w, r)
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

	payload, jsonErr := json.Marshal(params.request)

	if jsonErr != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "could not create a secure push notification")
	}

	userAPI.securePushNotification(w, string(payload), certificates, params.showCompleteOutput, params.showSplittedOutput)
}

func (userAPI *API) getSecurePushNotificationRequestParams(w http.ResponseWriter, r *http.Request) (*securePushNotificationRequestParams, bool) {

	// Decode request
	var securePNRequest model.SecurePushNotificationRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&securePNRequest); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	if !securePNRequest.IsValid() {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid secure push notification request")
		return nil, false
	}

	requestParams := extractUserParams(r, false)
	if requestParams == nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid requests params")
		return nil, false
	}

	showCompleteOutput := true
	showSplittedOutput := false

	if completeOutput := r.URL.Query().Get(QueryParamShowCompleteOutput); completeOutput == "false" {
		showCompleteOutput = false
	}

	if splittedOutput := r.URL.Query().Get(QueryParamShowSplittedOutput); splittedOutput == "true" {
		showSplittedOutput = true
	}

	return &securePushNotificationRequestParams{
		userId:             requestParams.userId,
		devices:            requestParams.devices,
		keyId:              requestParams.keyId,
		useAllDevices:      requestParams.useAllDevices,
		useDefaultDevice:   requestParams.useDefaultDevice,
		useDefaultKey:      requestParams.useDefaultKey,
		request:            securePNRequest,
		showCompleteOutput: showCompleteOutput,
		showSplittedOutput: showSplittedOutput,
	}, true
}

type securePushNotificationRequestParams struct {
	userId             string
	devices            []string
	keyId              string
	useAllDevices      bool
	useDefaultDevice   bool
	useDefaultKey      bool
	request            model.SecurePushNotificationRequest
	showCompleteOutput bool
	showSplittedOutput bool
}

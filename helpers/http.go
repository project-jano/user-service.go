package helpers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

const (
	PathParamUserId     = "userId"
	PathParamDeviceId   = "deviceId"
	QueryParamDeviceIds = "deviceIds"
	QueryParamKeyId     = "keyId"

	QueryValueAll = "all"
	DefaultKeyId  = "default"

	ContentType        = "Content-Type"
	DefaultContentType = "application/json; charset=UTF-8"
)

func ExtractUserId(r *http.Request) string {
	return extractPathParam(r, PathParamUserId)
}

func ExtractDeviceId(r *http.Request) string {
	return extractPathParam(r, PathParamDeviceId)
}

func ExtractDevicesFilters(r *http.Request) (bool, bool, []string, bool) {
	useAllDevices := false
	useDefaultDevice := true
	var devices []string

	if devicesString := r.URL.Query().Get(QueryParamDeviceIds); len(devicesString) > 0 {
		r := regexp.MustCompile(`[^\s,]+`)
		devices = r.FindAllString(devicesString, -1)
		useDefaultDevice = false
	}

	if ContainsStringInStringArray(devices, QueryValueAll) {
		useAllDevices = true
		useDefaultDevice = false
		devices = []string{}
	}

	return useAllDevices, useDefaultDevice, devices, true
}

func ExtractKeyId(r *http.Request) (bool, string) {
	useDefaultKey := true
	keyId := DefaultKeyId
	if keyIdString := r.URL.Query().Get(QueryParamKeyId); len(keyIdString) > 0 {
		if keyIdString != DefaultKeyId {
			keyId = keyIdString
			useDefaultKey = false
		}
	}
	return useDefaultKey, keyId
}

func extractPathParam(r *http.Request, pathParam string) string {
	params := mux.Vars(r)
	value, ok := params[pathParam]
	if !ok || strings.TrimSpace(value) == "" {
		return ""
	}
	decoded, error := url.QueryUnescape(value)
	if error != nil {
		return ""
	}
	return decoded
}

func SetupDefaultContentType(w http.ResponseWriter) {
	w.Header().Set(ContentType, DefaultContentType)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(map[string]string{"error_message": message, "error_id": uuid.NewString()})

	SetupDefaultContentType(w)
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func ResponseWithJSON(w http.ResponseWriter, jsonObject interface{}) {
	jsonErr := json.NewEncoder(w).Encode(jsonObject)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

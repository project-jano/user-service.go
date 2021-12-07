package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

func extractUserId(r *http.Request) string {
	return extractPathParam(r, QueryParamUserId)

}

func extractDeviceId(r *http.Request) string {
	return extractPathParam(r, QueryParamDeviceId)
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

func containsStringInStringArray(arr []string, key string) bool {
	for _, str := range arr {
		if str == key {
			return true
		}
	}
	return false
}

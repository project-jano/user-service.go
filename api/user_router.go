package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func (api *API) addUserRouter() {

	const version = "/v2"
	const path = "/users"

	userRouter := mux.NewRouter().PathPrefix(version).PathPrefix(path).Subrouter().StrictSlash(true)

	// Authenticated routes
	routes := Routes{
		Route{
			"SignUserCertificate",
			"POST",
			"/{userId}/devices/{deviceId}/csr",
			api.SignUserCertificate,
			true,
		},

		Route{
			"SecureMessageForUser",
			"POST",
			"/{userId}/messages/secure",
			api.SecureMessageForUser,
			true,
		},

		Route{
			"SecurePushNotificationForUser",
			"POST",
			"/{userId}/push-notifications/secure",
			api.SecurePushNotificationForUser,
			true,
		},

		Route{
			"DecodeSecureMessage",
			"POST",
			"/{userId}/messages/decode",
			api.DecodeSecureMessage,
			true,
		},
	}

	api.appendRoutes(userRouter, routes, true)

	api.Router.PathPrefix(version).PathPrefix(path).Handler(negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			if !api.Configuration.AuthEnabled {
				next(w, r)
				return
			}

			if api.basicAuth(w, r, api.Configuration.AuthUsername, api.Configuration.AuthPassword, path) {
				next(w, r)
			}
		}),
		negroni.Wrap(userRouter),
	))

}

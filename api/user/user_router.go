package user

import (
	"github.com/project-jano/user-service.go/repository"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/project-jano/user-service.go/app"
	"github.com/project-jano/user-service.go/auth"
	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/model"
	"github.com/urfave/negroni"
)

type API struct {
	fingerprint   string
	configuration app.APIConfiguration
	repository    repository.UserRepository
}

const (
	RSACipher         = "RSA"
	CertificatePrefix = "-----BEGIN CERTIFICATE-----"
	CertificateSufix  = "-----END CERTIFICATE-----"

	APIVersion = "/v2"
	APIPath = "/users"

	WWWAuthenticateHeader = "WWW-Authenticate"
)

func NewUserAPI(router *mux.Router, configuration app.APIConfiguration, fingerprint string, userRepository repository.UserRepository) {
	userAPI := API{
		fingerprint:   fingerprint,
		configuration: configuration,
		repository:    userRepository,
	}
	userAPI.addToRouter(router, configuration.TraceCallsEnabled)
}

func (userAPI *API) addToRouter(router *mux.Router, traceCalls bool) {

	userRouter := mux.NewRouter().PathPrefix(APIVersion).PathPrefix(APIPath).Subrouter().StrictSlash(true)

	// Authenticated routes
	routes := model.Routes{
		model.Route{
			Name:          "SignUserCertificate",
			Method:        "POST",
			Pattern:       "/{userId}/devices/{deviceId}/csr",
			HandlerFunc:   userAPI.SignUserCertificate,
			Authenticated: true,
		},

		model.Route{
			Name:          "SecureMessageForUser",
			Method:        "POST",
			Pattern:       "/{userId}/messages/secure",
			HandlerFunc:   userAPI.SecureMessageForUser,
			Authenticated: true,
		},

		model.Route{
			Name:          "SecurePushNotificationForUser",
			Method:        "POST",
			Pattern:       "/{userId}/push-notifications/secure",
			HandlerFunc:   userAPI.SecurePushNotificationForUser,
			Authenticated: true,
		},

		model.Route{
			Name:          "DecodeSecureMessage",
			Method:        "POST",
			Pattern:       "/{userId}/messages/decode",
			HandlerFunc:   userAPI.DecodeSecureMessage,
			Authenticated: true,
		},
	}

	helpers.AppendRoutes(userRouter, routes, false, traceCalls)

	router.PathPrefix(APIVersion).PathPrefix(APIPath).Handler(negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			if !userAPI.configuration.AuthEnabled {
				next(w, r)
				return
			}

			if !auth.IsAuthenticated(r, userAPI.configuration.AuthUsername, userAPI.configuration.AuthPassword, r.URL.Path) {
				w.Header().Set(WWWAuthenticateHeader, `Basic realm="`+ r.URL.Path +`, charset="UTF-8"`)
				helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorised")
				return
			}

			next(w, r)
		}),
		negroni.Wrap(userRouter),
	))

}

package api

import (
	"crypto/subtle"
	"net/http"
)


func (api * API)basicAuth(w http.ResponseWriter, r *http.Request, username, password, realm string) bool {

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`, charset="UTF-8"`)
		api.respondWithError(w, http.StatusUnauthorized, "Unauthorised")
		return false
	}

	return true
}
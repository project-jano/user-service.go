package auth

import (
	"crypto/subtle"
	"net/http"
)

func IsAuthenticated(r *http.Request, username, password, realm string) bool {

	user, pass, ok := r.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		return false
	}
	return true
}

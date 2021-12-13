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
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/project-jano/user-service.go/helpers"

	"github.com/project-jano/user-service.go/model"
	"github.com/project-jano/user-service.go/security"
)

func (userAPI *API) SignUserCertificate(w http.ResponseWriter, r *http.Request) {
	helpers.SetupDefaultContentType(w)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var csr model.CertificateSigningRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&csr); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	requestParams := extractUserParams(r, true)
	if requestParams == nil || len(requestParams.devices) != 1 {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid requests params")
		return
	}

	chainString, userCertificate := userAPI.signCertificateRequest(csr, requestParams)
	if userCertificate == nil || len(chainString) == 0 {
		helpers.RespondWithError(w, http.StatusInternalServerError, "failed to create Certificate")
		return
	}

	user := userAPI.repository.FindUser(ctx, requestParams.userId)
	if user == nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "failed to find User")
		return
	}

	user.Certificates = model.UserCertificatesAppend(user.Certificates, *userCertificate)

	updateError := userAPI.repository.UpdateCertificates(ctx, user)
	if updateError != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "failed to save Certificate")
		return
	}

	certificateSigningResponse := model.CertificateSigningResponse{
		Chain: chainString,
	}

	helpers.ResponseWithJSON(w, certificateSigningResponse)
}

func (userAPI *API) signCertificateRequest(csr model.CertificateSigningRequest, requestParams *userRequestParams) (string, *model.UserCertificate) {
	chain, err := security.SignCertificateRequest(csr, userAPI.configuration.ClientCertificateDuration, userAPI.configuration.Certificate, userAPI.configuration.PrivateKey)
	if err != nil {

		return "", nil
	}
	chainString := string(chain)

	var chainElements []string

	for _, element := range strings.Split(chainString, "\n-----END CERTIFICATE-----\n") {
		rawCert := strings.ReplaceAll(element, "-----BEGIN CERTIFICATE-----\n", "")
		if rawCert == "" {
			continue
		}
		chainElements = append(chainElements, rawCert)
	}

	userCertificate := &model.UserCertificate{
		KeyId:              csr.KeyId,
		DeviceId:           requestParams.devices[0],
		Cipher:             csr.Cipher,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		Certificate:        chainElements[0],
		Default_:           csr.Default_,
		Device:             csr.Device,
		Created:            time.Now(),
	}
	return chainString, userCertificate
}

package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *API) SecureMessageForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userId := extractUserId(r)
	if userId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid userId")
		return
	}

	// Decode request
	var secureMessageRequest model.SecureMessageRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&secureMessageRequest); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if !secureMessageRequest.IsValid() {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Get devices in requests. If empty means all
	allDevices := false
	var devices []string
	if devicesString := r.URL.Query().Get("devices"); len(devicesString) > 0 {
		r := regexp.MustCompile(`[^\s,]+`)
		devices = r.FindAllString(devicesString, -1)
	}

	if containsStringInStringArray(devices, "all") {
		allDevices = true
		devices = []string{}
	}

	ctx := context.TODO()
	var user model.User
	filter := bson.D{{Key: "userId", Value: userId}}

	findError := a.UserDatabase.FindOne(ctx, filter).Decode(&user)
	if findError != nil {
		if findError != mongo.ErrNoDocuments {
			a.respondWithJSON(w, http.StatusNotFound, "user not found")
			return
		}
		a.respondWithJSON(w, http.StatusInternalServerError, "failed to fetch User")
		return
	}

	if len(user.Certificates) == 0 {
		a.respondWithJSON(w, http.StatusNotFound, "no certificates found")
		return
	}

	// Get all certs for those devices
	var certificates []model.UserCertificate
	for _, cert := range user.Certificates {
		var appendCert = false
		if allDevices {
			appendCert = true
		} else if len(devices) == 0 && cert.Default_ {
			appendCert = true
		} else if containsStringInStringArray(devices, cert.DeviceId) {
			appendCert = true
		}

		if appendCert {
			certificates = append(certificates, cert)
		}
	}

	if len(certificates) == 0 {
		a.respondWithJSON(w, http.StatusBadRequest, "no devices found")
		return
	}

	payloadToEncrypt := model.Payload{
		Message:     secureMessageRequest.Message,
		Timestamp:   time.Now().Unix(),
		Fingerprint: a.Fingerprint,
	}

	jsonPayloadToEncrypt, jsonErr := json.Marshal(&payloadToEncrypt)
	if jsonErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "cannot create payload to encrypt")
		return
	}

	var securedPayloads []model.SecuredPayload
	for _, userCert := range certificates {

		pemStr := "-----BEGIN CERTIFICATE-----\n" + userCert.Certificate + "\n-----END CERTIFICATE-----"

		block, _ := pem.Decode([]byte(pemStr))
		if block == nil {
			a.respondWithJSON(w, http.StatusInternalServerError, "could not decode PEM")
			return
		}
		var cert *x509.Certificate
		cert, _ = x509.ParseCertificate(block.Bytes)

		if strings.HasPrefix(userCert.Cipher, "RSA") {
			publicKey := cert.PublicKey.(*rsa.PublicKey)

			encrypted, encryptedErr := rsa.EncryptPKCS1v15(rand.Reader, publicKey, jsonPayloadToEncrypt)
			if encryptedErr != nil {
				a.respondWithJSON(w, http.StatusInternalServerError, fmt.Sprintf("could not encrypt. %+v", encryptedErr))
				return
			}

			hashAlg, hashed := helpers.HashBytesUsingAlgorithm(encrypted, userCert.SignatureAlgorithm)

			signature, err := rsa.SignPKCS1v15(rand.Reader, a.Configuration.PrivateKey, hashAlg, hashed[:])
			if err != nil {
				a.respondWithJSON(w, http.StatusInternalServerError, fmt.Sprintf("could not sign. %+v", encryptedErr))
				return
			}

			securedPayload := model.SecuredPayload{
				KeyId:     userCert.KeyId,
				DeviceId:  userCert.DeviceId,
				Payload:   base64.StdEncoding.EncodeToString(encrypted),
				Signature: base64.StdEncoding.EncodeToString(signature),
			}

			securedPayloads = append(securedPayloads, securedPayload)
		}
	}

	responseJsonErr := json.NewEncoder(w).Encode(securedPayloads)
	if responseJsonErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "failed to create response")
	}
}

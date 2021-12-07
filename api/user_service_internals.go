package api

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

func (a *API) securePayload(w http.ResponseWriter, r *http.Request, payloadString string, certificates []model.UserCertificate) {

	payloadToEncrypt := model.Payload{
		Message:     payloadString,
		Timestamp:   time.Now().Unix(),
		Fingerprint: a.Fingerprint,
	}

	jsonPayloadToEncrypt, jsonErr := json.Marshal(&payloadToEncrypt)
	if jsonErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "cannot create payload to encrypt")
		return
	}

	var securedPayloads []model.SecuredMessage
	for _, userCert := range certificates {

		pemStr := CertificatePrefix + "\n" + userCert.Certificate + "\n" + CertificateSufix

		block, _ := pem.Decode([]byte(pemStr))
		if block == nil {
			a.respondWithJSON(w, http.StatusInternalServerError, "could not decode PEM")
			return
		}
		var cert *x509.Certificate
		cert, _ = x509.ParseCertificate(block.Bytes)

		if strings.HasPrefix(userCert.Cipher, RSACipher) {
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

			securedPayload := model.SecuredMessage{
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

func (a *API) securePushNotification(w http.ResponseWriter, r *http.Request, payloadString string, certificates []model.UserCertificate, showCompleteOutput bool, showSplittedOutput bool) {

	payloadToEncrypt := model.Payload{
		Message:     payloadString,
		Timestamp:   time.Now().Unix(),
		Fingerprint: a.Fingerprint,
	}

	jsonPayloadToEncrypt, jsonErr := json.Marshal(&payloadToEncrypt)
	if jsonErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "cannot create payload to encrypt")
		return
	}

	var securedPushNotifications []model.SecuredPushNotification
	const maxChunkSize = 200

	for _, userCert := range certificates {

		pemStr := CertificatePrefix + "\n" + userCert.Certificate + "\n" + CertificateSufix

		block, _ := pem.Decode([]byte(pemStr))
		if block == nil {
			a.respondWithJSON(w, http.StatusInternalServerError, "could not decode PEM")
			return
		}
		var cert *x509.Certificate
		cert, _ = x509.ParseCertificate(block.Bytes)

		if strings.HasPrefix(userCert.Cipher, RSACipher) {
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

			base64Encrypted := base64.StdEncoding.EncodeToString(encrypted)
			base64Signature := base64.StdEncoding.EncodeToString(signature)

			p := model.CompleteAndSplittedValue{
				Complete: base64Encrypted,
			}

			s := model.CompleteAndSplittedValue{
				Complete: base64Signature,
			}

			if showSplittedOutput {
				p.Splitted = helpers.SplitString(base64Encrypted, maxChunkSize)
				s.Splitted = helpers.SplitString(base64Signature, maxChunkSize)
			}
			if !showCompleteOutput {
				p.Complete = ""
				s.Complete = ""
			}

			securedPayload := model.SecuredPushNotification{
				KeyId:     userCert.KeyId,
				DeviceId:  userCert.DeviceId,
				Payload:   p,
				Signature: s,
			}

			securedPushNotifications = append(securedPushNotifications, securedPayload)
		}
	}

	responseJsonErr := json.NewEncoder(w).Encode(securedPushNotifications)
	if responseJsonErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "failed to create response")
	}
}

func (a *API) findUser(userId string, w http.ResponseWriter) (*model.User, bool) {
	ctx := context.TODO()
	var user model.User
	filter := bson.D{{Key: "userId", Value: userId}}

	findError := a.UserDatabase.FindOne(ctx, filter).Decode(&user)
	if findError != nil {
		if findError != mongo.ErrNoDocuments {
			a.respondWithJSON(w, http.StatusNotFound, "user not found")
			return nil, false
		}
		a.respondWithJSON(w, http.StatusInternalServerError, "failed to fetch User")
		return nil, false
	}

	if len(user.Certificates) == 0 {
		a.respondWithJSON(w, http.StatusNotFound, "no certificates found")
		return nil, false
	}

	return &user, true
}

func (a *API) filterCertificates(allCertificates []model.UserCertificate, allDevices bool, useDefaultKey bool, devices []string, keyId string) []model.UserCertificate {
	var certificates []model.UserCertificate

	for _, cert := range allCertificates {
		var deviceMatchesCriteria = false

		if allDevices {
			deviceMatchesCriteria = true
		} else if len(devices) == 0 && cert.Default_ {
			deviceMatchesCriteria = true
		} else if containsStringInStringArray(devices, cert.DeviceId) {
			deviceMatchesCriteria = true
		}

		if !deviceMatchesCriteria {
			continue
		}

		if useDefaultKey {
			if keyId == DefaultKeyId {
				certificates = append(certificates, cert)
			}
			continue
		}

		if keyId == cert.KeyId {
			certificates = append(certificates, cert)
		}
	}
	return certificates
}

func (a *API) getSecureMessageRequestParams(w http.ResponseWriter, r *http.Request) (*secureMessageRequestParams, bool) {

	userId := extractUserId(r)
	if userId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid userId")
		return nil, false
	}

	// Decode request
	var secureMessageRequest model.SecureMessageRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&secureMessageRequest); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	if !secureMessageRequest.IsValid() {
		a.respondWithError(w, http.StatusBadRequest, "invalid secure message request")
		return nil, false
	}

	// Get devices in requests. If empty means all
	useAllDevices := false
	useDefaultDevice := true

	var devices []string
	if devicesString := r.URL.Query().Get(QueryParamDeviceIds); len(devicesString) > 0 {
		r := regexp.MustCompile(`[^\s,]+`)
		devices = r.FindAllString(devicesString, -1)
		useDefaultDevice = false
	}

	if containsStringInStringArray(devices, QueryValueAll) {
		useAllDevices = true
		useDefaultDevice = false
		devices = []string{}
	}

	if !useDefaultDevice && !useAllDevices && len(devices) == 0 {
		a.respondWithError(w, http.StatusBadRequest, "no device matches deviceId criteria")
		return nil, false
	}

	// Get keys in requests. If empty means default
	useDefaultKey := true
	keyId := DefaultKeyId
	if keyIdString := r.URL.Query().Get(QueryParamKeyId); len(keyIdString) > 0 {
		if keyIdString != DefaultKeyId {
			keyId = keyIdString
			useDefaultKey = false
		}
	}

	return &secureMessageRequestParams{
		userId:           userId,
		devices:          devices,
		keyId:            keyId,
		useAllDevices:    useAllDevices,
		useDefaultDevice: useDefaultDevice,
		useDefaultKey:    useDefaultKey,
		request:          secureMessageRequest,
	}, true
}

func (a *API) getSecurePushNotificationRequestParams(w http.ResponseWriter, r *http.Request) (*securePushNotificationRequestParams, bool) {

	userId := extractUserId(r)
	if userId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid userId")
		return nil, false
	}

	// Decode request
	var securePNRequest model.SecurePushNotificationRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&securePNRequest); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return nil, false
	}

	if !securePNRequest.IsValid() {
		a.respondWithError(w, http.StatusBadRequest, "invalid secure push notification request")
		return nil, false
	}

	// Get devices in requests. If empty means all
	useAllDevices := false
	useDefaultDevice := true

	var devices []string
	if devicesString := r.URL.Query().Get(QueryParamDeviceIds); len(devicesString) > 0 {
		r := regexp.MustCompile(`[^\s,]+`)
		devices = r.FindAllString(devicesString, -1)
		useDefaultDevice = false
	}

	if containsStringInStringArray(devices, QueryValueAll) {
		useAllDevices = true
		useDefaultDevice = false
		devices = []string{}
	}

	if !useDefaultDevice && !useAllDevices && len(devices) == 0 {
		a.respondWithError(w, http.StatusBadRequest, "no device matches deviceId criteria")
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

	// Get keys in requests. If empty means default
	useDefaultKey := true
	keyId := DefaultKeyId
	if keyIdString := r.URL.Query().Get(QueryParamKeyId); len(keyIdString) > 0 {
		if keyIdString != DefaultKeyId {
			keyId = keyIdString
			useDefaultKey = false
		}
	}

	return &securePushNotificationRequestParams{
		userId:             userId,
		devices:            devices,
		keyId:              keyId,
		useAllDevices:      useAllDevices,
		useDefaultDevice:   useDefaultDevice,
		useDefaultKey:      useDefaultKey,
		request:            securePNRequest,
		showCompleteOutput: showCompleteOutput,
		showSplittedOutput: showSplittedOutput,
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

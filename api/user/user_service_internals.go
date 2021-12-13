package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/model"
)

func (userAPI *API) securePayload(w http.ResponseWriter, payloadString string, certificates []model.UserCertificate) {

	payloadToEncrypt := model.Payload{
		Message:     payloadString,
		Timestamp:   time.Now().Unix(),
		Fingerprint: userAPI.fingerprint,
	}

	jsonPayloadToEncrypt, jsonErr := json.Marshal(&payloadToEncrypt)
	if jsonErr != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot create payload to encrypt")
		return
	}

	var securedPayloads []model.SecuredMessage
	for _, userCert := range certificates {

		cert, ok := userCertificateToX509(w, userCert)
		if !ok {
			return
		}

		if strings.HasPrefix(userCert.Cipher, RSACipher) {
			publicKey := cert.PublicKey.(*rsa.PublicKey)

			encrypted, signature, ok := userAPI.encryptAndSign(w, publicKey, jsonPayloadToEncrypt, userCert)
			if !ok {
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

	helpers.ResponseWithJSON(w, securedPayloads)
}

func (userAPI *API) securePushNotification(w http.ResponseWriter, payloadString string, certificates []model.UserCertificate, showCompleteOutput bool, showSplittedOutput bool) {

	payloadToEncrypt := model.Payload{
		Message:     payloadString,
		Timestamp:   time.Now().Unix(),
		Fingerprint: userAPI.fingerprint,
	}

	jsonPayloadToEncrypt, jsonErr := json.Marshal(&payloadToEncrypt)
	if jsonErr != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "cannot create payload to encrypt")
		return
	}

	var securedPushNotifications []model.SecuredPushNotification
	const maxChunkSize = 200

	for _, userCert := range certificates {

		var cert *x509.Certificate
		var certOK bool

		if cert, certOK = userCertificateToX509(w, userCert); !certOK {
			continue
		}

		if strings.HasPrefix(userCert.Cipher, RSACipher) {
			publicKey := cert.PublicKey.(*rsa.PublicKey)

			var encrypted []byte
			var signature []byte
			var encryptOK bool

			if encrypted, signature, encryptOK = userAPI.encryptAndSign(w, publicKey, jsonPayloadToEncrypt, userCert); !encryptOK {
				continue
			}

			base64Encrypted := base64.StdEncoding.EncodeToString(encrypted)
			base64Signature := base64.StdEncoding.EncodeToString(signature)

			securedPayload := createSecuredPushNotification(base64Encrypted, base64Signature, showSplittedOutput, maxChunkSize, showCompleteOutput, userCert)

			securedPushNotifications = append(securedPushNotifications, securedPayload)
		}

	}

	helpers.ResponseWithJSON(w, securedPushNotifications)
}

func (userAPI *API) encryptAndSign(w http.ResponseWriter, publicKey *rsa.PublicKey, jsonPayloadToEncrypt []byte, userCert model.UserCertificate) ([]byte, []byte, bool) {
	encrypted, encryptedErr := rsa.EncryptPKCS1v15(rand.Reader, publicKey, jsonPayloadToEncrypt)
	if encryptedErr != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not encrypt. %+v", encryptedErr))
		return nil, nil, false
	}

	hashAlg, hashed := helpers.HashBytesUsingAlgorithm(encrypted, userCert.SignatureAlgorithm)

	signature, err := rsa.SignPKCS1v15(rand.Reader, userAPI.configuration.PrivateKey, hashAlg, hashed[:])
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not sign. %+v", encryptedErr))
		return nil, nil, false
	}
	return encrypted, signature, true
}

func filterCertificates(allCertificates []model.UserCertificate, allDevices bool, useDefaultKey bool, devices []string, keyId string) []model.UserCertificate {
	var certificates []model.UserCertificate

	for _, cert := range allCertificates {
		var deviceMatchesCriteria = false

		if allDevices {
			deviceMatchesCriteria = true
		} else if len(devices) == 0 && cert.Default_ {
			deviceMatchesCriteria = true
		} else if helpers.ContainsStringInStringArray(devices, cert.DeviceId) {
			deviceMatchesCriteria = true
		}

		if !deviceMatchesCriteria {
			continue
		}

		if useDefaultKey {
			if keyId == helpers.DefaultKeyId {
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

func userCertificateToX509(w http.ResponseWriter, userCert model.UserCertificate) (*x509.Certificate, bool) {
	pemStr := CertificatePrefix + "\n" + userCert.Certificate + "\n" + CertificateSufix

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "could not decode PEM")
		return nil, false
	}
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	return cert, true
}

func createSecuredPushNotification(base64Encrypted string, base64Signature string, showSplittedOutput bool, maxChunkSize int, showCompleteOutput bool, userCert model.UserCertificate) model.SecuredPushNotification {
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
	return securedPayload
}

func extractUserParams(r *http.Request, allowOnlyOneDevice bool) *userRequestParams {

	userId := helpers.ExtractUserId(r)
	if userId == "" {
		return nil
	}

	// Get devices in requests. If empty means all
	useAllDevices, useDefaultDevice, devices, ok := helpers.ExtractDevicesFilters(r)
	if !ok || allowOnlyOneDevice {
		uniqueDevice := helpers.ExtractDeviceId(r)
		if uniqueDevice == "" {
			return nil
		}
		useAllDevices = false
		useDefaultDevice = false
		devices = []string{uniqueDevice}
	}

	// Get keys in requests. If empty means default
	useDefaultKey, keyId := helpers.ExtractKeyId(r)

	return &userRequestParams{
		userId:           userId,
		devices:          devices,
		useAllDevices:    useAllDevices,
		useDefaultDevice: useDefaultDevice,
		keyId:            keyId,
		useDefaultKey:    useDefaultKey,
	}
}

type userRequestParams struct {
	userId           string
	devices          []string
	useAllDevices    bool
	useDefaultDevice bool
	keyId            string
	useDefaultKey    bool
}

package model

import (
	"fmt"
	"time"
)

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type UserCertificate struct {

	// Identifier of this certificate in the user's device
	KeyId string `json:"keyId" bson:"keyId"`
	// Identifier of this certificate in the user's device
	DeviceId string `json:"deviceId" bson:"deviceId"`

	// Identifier if the transformation that should be used in the cipher
	Cipher string `json:"cipher" bson:"cipher"`

	// Identifier of the Signature Algorithm
	SignatureAlgorithm string `json:"signatureAlgorithm" bson:"signatureAlgorithm"`

	// User's certificate
	Certificate string `json:"certificate" bson:"certificate"`

	// true if this is the default credential for this user.
	Default_ bool `json:"default,omitempty" bson:"default"`

	// Information about the device
	Device *Device `json:"device,omitempty" bson:"device"`

	// When the entry was created
	Created time.Time `json:"created,omitempty" bson:"created"`
}

func UserCertificatesAppend(arr []UserCertificate, newCert UserCertificate) []UserCertificate {
	occurred := map[string]bool{}
	var result []UserCertificate

	newCertKey := fmt.Sprintf("key=%s;device=%s", newCert.KeyId, newCert.DeviceId)

	for _, cert := range arr {
		certKey := fmt.Sprintf("key=%s;device=%s", cert.KeyId, cert.DeviceId)

		// Keep only one cert as default
		if cert.Default_ && newCert.Default_ {
			cert.Default_ = false
		}

		// If 'cert' has a different keyId and deviceId as the 'newCert' add it to the list
		if certKey != newCertKey && !occurred[certKey] {
			occurred[certKey] = true
			result = append(result, cert)
		}
	}

	// certs will always be appended to this list in order of creation
	result = append(result, newCert)

	return result
}

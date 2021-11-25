package model

import (
	"fmt"
	"time"
)

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
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
	newCertKey := fmt.Sprintf("key=%s;device=%s", newCert.KeyId, newCert.DeviceId)
	var result []UserCertificate
	for _, cert := range arr {
		certKey := fmt.Sprintf("key=%s;device=%s", cert.KeyId, cert.DeviceId)
		if certKey != newCertKey && occurred[certKey] == false {
			occurred[certKey] = true
			result = append(result, cert)
		}
	}
	result = append(result, newCert)
	return result
}

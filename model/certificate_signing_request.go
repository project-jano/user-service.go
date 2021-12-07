package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

type CertificateSigningRequest struct {
	// Identifier of this certificate in the user's device
	KeyId string `json:"keyId"`

	// Identifier if the transformation that should be used in the cipher
	Cipher string `json:"cipher"`

	// Identifier of the signature algorithm used when signing a ciphered message
	SignatureAlgorithm string `json:"signatureAlgorithm"`

	// Certificate Signing Request
	Request string `json:"request"`

	// Information about the device
	Device *Device `json:"device,omitempty"`

	// true if this is the default credential for this user.
	Default_ bool `json:"default,omitempty"`
}

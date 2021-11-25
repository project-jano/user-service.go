package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type SecuredPayload struct {
	// Identifier of this certificate in the user's device
	KeyId string `json:"keyId"`

	// Identifier of user's device
	DeviceId string `json:"deviceId"`

	// Payload is a ciphered string representation of Payload struct
	Payload string `json:"payload"`

	// Signature of the plain text payload
	Signature string `json:"signature"`
}

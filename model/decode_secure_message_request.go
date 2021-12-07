package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type DecodeSecuredMessageRequest struct {
	// Identifier of this certificate in the user's device
	KeyId string `json:"keyId"`

	// SecuredMessage is a ciphered string representation of Payload struct
	SecuredMessage string `json:"securedMessage"`

	// Signature of the plain text payload
	Signature string `json:"signature"`
}

func (request *DecodeSecuredMessageRequest) IsValid() bool {
	return len(request.KeyId) > 0 &&
		len(request.SecuredMessage) > 0 &&
		len(request.Signature) > 0
}

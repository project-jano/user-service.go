package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type DecodeSecurePayloadRequest struct {
	// Identifier of this certificate in the user's device
	KeyId string `json:"keyId"`

	// SecuredPayload is a ciphered string representation of Payload struct
	SecuredPayload string `json:"securedPayload"`

	// Signature of the plain text payload
	Signature string `json:"signature"`
}

func (request *DecodeSecurePayloadRequest) IsValid() bool {
	return len(request.KeyId) > 0 &&
		len(request.SecuredPayload) > 0 &&
		len(request.Signature) > 0
}

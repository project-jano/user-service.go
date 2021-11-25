package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type SecureMessageRequest struct {
	// message to secure. Either plain text, JSON string, or any other type of structure represented as a string.
	Message string `json:"message"`
}

func (request *SecureMessageRequest) IsValid() bool {
	return len(request.Message) > 0
}

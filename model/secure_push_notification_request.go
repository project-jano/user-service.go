package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type SecurePushNotificationRequest struct {
	//Title of the Push Notification
	Title string `json:"title,omitempty"`

	// Body of the Push Notification
	Body string `json:"body,omitempty"`

	// A JSON (string) of parameters
	Payload string `json:"payload,omitempty"`
}

func (request *SecurePushNotificationRequest) IsValid() bool {
	return len(request.Title) > 0 || len(request.Body) > 0 || len(request.Payload) > 0
}

package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type LivenessResponse struct {
	// Status of the service. Always "up" when the response is OK
	Status string `json:"status"`

	// Name of the host where the service run
	Hostname string `json:"hostname"`
}

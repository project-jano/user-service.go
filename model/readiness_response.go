package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type ReadinessResponse struct {
	// Status of the service. Always "ready" when the response is OK. Ready means the service is up and running, connected to the Database, and ready to receive requests.
	Status string `json:"status"`

	// Name of the host where the service run
	Hostname string `json:"hostname"`
}

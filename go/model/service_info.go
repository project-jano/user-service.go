package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

// ServiceInfo includes a string with the messageand two properties, Server TimeStamp and Fingerprint, in order to generate entrophy in the Payload string
type ServiceInfo struct {

	// Unix timestamp of the server when securing the payload in mS.
	Timestamp int64 `json:"timestamp"`

	// Hostname where the service runs
	Hostname string `json:"hostname"`

	// Fingerprint of the service where the Payload was encrypted
	Fingerprint string `json:"fingerprint"`
}

package model

/* Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

// Payload includes a string with the message, Server TimeStamp and Fingerprint, in order to generate entrophy in the Payload string
type Payload struct {

	// Unix timestamp of the server when securing the payload in mS.
	Timestamp int64 `json:"timestamp"`

	// Fingerprint of the service where the Payload was encrypted
	Fingerprint string `json:"fingerprint"`

	// Encrypted payload
	Message string `json:"message"`
}

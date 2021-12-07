package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

type CertificateSigningResponse struct {
	// Chain with user certificate and certificate used for signing
	Chain string `json:"chain,omitempty"`
}

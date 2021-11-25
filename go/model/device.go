package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type Device struct {
	// Platform of the device where the certificate was created
	Platform string `json:"platform" bson:"platform"`
	// User agent of the device where the certificate was generated
	UserAgent string `json:"userAgent" bson:"userAgent"`
}

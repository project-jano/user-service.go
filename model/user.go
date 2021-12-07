package model

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

type User struct {
	UserId string `json:"userId" bson:"userId"`

	Certificates []UserCertificate `json:"certificates" bson:"certificates"`
}

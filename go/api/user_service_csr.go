package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/project-jano/user-service.go/go/model"
	"github.com/project-jano/user-service.go/go/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *API) SignUserCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userId := extractUserId(r)
	if userId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid userId")
		return
	}
	deviceId := extractDeviceId(r)
	if deviceId == "" {
		a.respondWithError(w, http.StatusBadRequest, "invalid deviceId")
		return
	}

	var csr model.CertificateSigningRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&csr); err != nil {
		a.respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	chain, err := security.SignCertificateRequest(csr, a.Configuration.ClientCertificateDuration, a.Configuration.Certificate, a.Configuration.PrivateKey)
	if err != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to create Certificate. %v", err))
		return
	}
	chainString := string(chain)

	var chainElements []string

	for _, element := range strings.Split(chainString, "\n-----END CERTIFICATE-----\n") {
		rawCert := strings.ReplaceAll(element, "-----BEGIN CERTIFICATE-----\n", "")
		if rawCert == "" {
			continue
		}
		chainElements = append(chainElements, rawCert)
	}

	userCertificate := model.UserCertificate{
		KeyId:              csr.KeyId,
		DeviceId:           deviceId,
		Cipher:             csr.Cipher,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		Certificate:        chainElements[0],
		Default_:           csr.Default_,
		Device:             csr.Device,
		Created:            time.Now(),
	}

	var user *model.User
	filter := bson.D{{Key: "userId", Value: userId}}

	findError := a.UserDatabase.FindOne(context.TODO(), filter).Decode(&user)
	if findError != nil {
		if findError != mongo.ErrNoDocuments {
			a.respondWithJSON(w, http.StatusInternalServerError, "failed to find User")
			return
		}

		user = &model.User{
			UserId:       userId,
			Certificates: []model.UserCertificate{},
		}
	}

	user.Certificates = model.UserCertificatesAppend(user.Certificates, userCertificate)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "certificates", Value: user.Certificates}}}}
	opts := options.Update().SetUpsert(true)
	_, upsertErr := a.UserDatabase.UpdateOne(context.TODO(), filter, update, opts)
	if upsertErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "failed to save Certificate")
		return
	}

	jsonErr := json.NewEncoder(w).Encode(map[string]string{"chain": chainString})
	if jsonErr != nil {
		a.respondWithJSON(w, http.StatusInternalServerError, "failed to create response")
	}
}

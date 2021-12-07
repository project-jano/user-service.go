package api

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
)

type Configuration struct {
	AuthUsername string
	AuthPassword string
	AuthEnabled  bool

	TraceCallsEnabled bool

	CertificatePEM            string
	Certificate               *x509.Certificate
	PrivateKey                *rsa.PrivateKey
	ClientCertificateDuration time.Duration

	DatabaseName string
}

type API struct {
	Configuration Configuration
	Router        *mux.Router
	MongoClient   *mongo.Client
	UserDatabase  *mongo.Collection
	Fingerprint   string
}

func NewAPI(config Configuration, client *mongo.Client) *API {

	fingerprint := fmt.Sprintf("%x", CreateFingerprint())

	api := &API{
		Configuration: config,
		Router:        mux.NewRouter().StrictSlash(true),
		MongoClient:   client,
		UserDatabase:  client.Database(config.DatabaseName).Collection("users"),
		Fingerprint:   fingerprint,
	}

	api.addSecurityRouter()
	api.addUserRouter()
	api.addHealthRouter()

	// Endpoints where prometheus Middleware doesnot apply
	api.Router.Handle("/metrics", promhttp.Handler())

	api.Router.NotFoundHandler = Logger(http.NotFoundHandler(), "NotFound")

	return api
}

func (a *API) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *API) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set(ContentType, DefaultContentType)

	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func CreateFingerprint() [32]byte {
	hostname, _ := os.Hostname()
	goos := runtime.GOOS
	ip := ""

	addrs, _ := net.LookupIP(hostname)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip = ipv4.String()
		}
	}

	fingerprint := fmt.Sprintf("%s@%s%s", hostname, goos, ip)
	return sha256.Sum256([]byte(fingerprint))
}

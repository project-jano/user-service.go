package app

import (
	"crypto/rsa"
	"crypto/x509"
	"time"
)

type APIConfiguration struct {
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

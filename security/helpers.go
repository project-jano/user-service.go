// Package security provides primitives and helper functions for dealing with PKI (certificates, CSR, ciphering, etc)
package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/project-jano/user-service.go/model"
)

// SignCertificateRequest creates a certificate chain with the CSR signed with a root certificate
func SignCertificateRequest(csr model.CertificateSigningRequest, certificateDuration time.Duration, rootCertificate *x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	certificateRequest := parseCertificateRequest(csr.Request)

	if certificateRequest == nil {
		return nil, errors.New("invalid certificate request")
	}

	err := certificateRequest.CheckSignature()
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	notAfter := time.Now().Add(certificateDuration)

	random := rand.Reader

	// Create serial number for X.509 certificate
	serialNumber, err := rand.Int(random, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:       serialNumber,
		NotBefore:          notBefore,
		NotAfter:           notAfter,
		Subject:            certificateRequest.Subject,
		PublicKeyAlgorithm: certificateRequest.PublicKeyAlgorithm,
		PublicKey:          certificateRequest.PublicKey,
		SignatureAlgorithm: x509.SHA512WithRSA,
		DNSNames:           certificateRequest.DNSNames,
		IPAddresses:        certificateRequest.IPAddresses,
		URIs:               certificateRequest.URIs,
		EmailAddresses:     certificateRequest.EmailAddresses,
		ExtraExtensions:    certificateRequest.Extensions,
		KeyUsage:           x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
	}

	signedCsr, err := x509.CreateCertificate(rand.Reader, &template, rootCertificate, certificateRequest.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return createPEMChain([]*x509.Certificate{rootCertificate}, signedCsr)
}

// LoadCertificateAndKey creates x509 and rsa private key from  its PEM files. It validates that the public key of the two public keys matches.
func LoadCertificateAndKey(certificatePEM, privateKeyPEM string) (*x509.Certificate, *rsa.PrivateKey, error) {
	decodedCert, _ := pem.Decode([]byte(certificatePEM))
	decodedPrivateKey, _ := pem.Decode([]byte(privateKeyPEM))

	certificate, certificateErr := x509.ParseCertificate(decodedCert.Bytes)
	if certificateErr != nil {
		return nil, nil, fmt.Errorf("unable to load certificate. %v", certificateErr)
	}

	privateKey, privateKeyErr := x509.ParsePKCS1PrivateKey(decodedPrivateKey.Bytes)
	if privateKeyErr != nil {
		return nil, nil, fmt.Errorf("unable to load private key. %v", privateKeyErr)
	}

	certificatePublicKey := (certificate.PublicKey).(*rsa.PublicKey)
	if !certificatePublicKey.Equal(privateKey.Public()) {
		return nil, nil, errors.New("certificate and private key does not match")
	}

	return certificate, privateKey, nil
}

func parseCertificateRequest(csrPEM string) *x509.CertificateRequest {

	decoded, _ := pem.Decode([]byte(csrPEM))
	var csr *x509.CertificateRequest
	var err error
	if decoded != nil {
		csr, err = x509.ParseCertificateRequest(decoded.Bytes)
	} else {
		csr, err = x509.ParseCertificateRequest([]byte(csrPEM))
	}

	if err != nil {
		log.Print(err)
	}

	return csr
}

func createPEMChain(authorities []*x509.Certificate, cert []byte) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 2048))

	// encode cert
	err := pem.Encode(buf, &pem.Block{
		Type: "CERTIFICATE", Bytes: cert,
	})
	if err != nil {
		return nil, err
	}
	// encode intermediates
	for _, ca := range authorities {
		err := pem.Encode(buf, &pem.Block{
			Type: "CERTIFICATE", Bytes: ca.Raw,
		})
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

package pub

import (
	"github.com/project-jano/user-service.go/helpers"
	"os"
	"testing"
	"time"
)

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

func TestCreateServiceInfo(t *testing.T) {
	hostname, _ := os.Hostname()
	fingerprint := helpers.HashedFingerprint()

	serviceInfo, ok := createServiceInfo(fingerprint)
	if !ok {
		t.Fail()
	}

	if serviceInfo.Fingerprint != fingerprint {
		t.Fail()
	}
	if serviceInfo.Hostname != hostname {
		t.Fail()
	}
	if serviceInfo.Timestamp < time.Now().Unix() {
		t.Fail()
	}
}

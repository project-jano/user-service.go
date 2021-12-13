package helpers

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestCreateFingerprint(t *testing.T) {
	fingerprint := createFingerprint()
	if len(fingerprint) == 0 {
		t.Fail()
	}

	goos := runtime.GOOS
	hostname, _ := os.Hostname()

	if !strings.HasPrefix(fingerprint, goos) {
		t.Fail()
	}
	if !strings.Contains(fingerprint, fmt.Sprintf(fingerprintFormat, goos, hostname, "")) {
		t.Fail()
	}
}

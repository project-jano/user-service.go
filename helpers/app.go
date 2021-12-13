package helpers

import (
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"runtime"
)

const (
	fingerprintFormat = "%s:%s@%s"
)

func createFingerprint() string {
	goos := runtime.GOOS
	hostname, _ := os.Hostname()
	ip := ""

	addrs, _ := net.LookupIP(hostname)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip = ipv4.String()
		}
	}

	fingerprint := fmt.Sprintf(fingerprintFormat, goos, hostname, ip)
	return fingerprint
}

func HashedFingerprint() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(createFingerprint())))
}

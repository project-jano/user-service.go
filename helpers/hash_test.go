package helpers

import (
	"crypto"
	"testing"
)

func TestCreateSHA256(t *testing.T) {
	if hashAlgorithmWithIdentifier(SHA256) != crypto.SHA256 {
		t.Fail()
	}
}

func TestCreateSHA384(t *testing.T) {
	if hashAlgorithmWithIdentifier(SHA384) != crypto.SHA384 {
		t.Fail()
	}
}

func TestCreateSHA512(t *testing.T) {
	if hashAlgorithmWithIdentifier(SHA512) != crypto.SHA512 {
		t.Fail()
	}
}

func TestMissingSHA1(t *testing.T) {
	if hashAlgorithmWithIdentifier("SHA1") == crypto.SHA1 {
		t.Fail()
	}
}
